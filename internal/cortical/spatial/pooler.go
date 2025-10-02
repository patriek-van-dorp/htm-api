package spatial

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/htm-project/neural-api/internal/cortical/sdr"
	"github.com/htm-project/neural-api/internal/domain/htm"
	"gonum.org/v1/gonum/mat"
)

// SpatialPooler implements the HTM spatial pooler algorithm
type SpatialPooler struct {
	// Configuration
	config *htm.SpatialPoolerConfig

	// Core matrices (using gonum for matrix operations)
	connections       *mat.Dense // [columnCount x inputWidth] - synaptic connections
	permanences       *mat.Dense // [columnCount x inputWidth] - permanence values
	connectedSynapses *mat.Dense // [columnCount x inputWidth] - binary connected matrix

	// Column state
	activeDutyCycles     []float64 // [columnCount] - duty cycles for active columns
	overlapDutyCycles    []float64 // [columnCount] - duty cycles for overlap
	minOverlapThresholds []int     // [columnCount] - minimum overlap thresholds per column
	boostFactors         []float64 // [columnCount] - boosting factors

	// Learning state
	iterationNum   int64
	lastUpdateTime time.Time

	// Random number generator (for randomized mode)
	rng *rand.Rand

	// Performance metrics
	metrics *htm.SpatialPoolerMetrics
}

// NewSpatialPooler creates a new spatial pooler with the given configuration
func NewSpatialPooler(config *htm.SpatialPoolerConfig) (*SpatialPooler, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	sp := &SpatialPooler{
		config:               config,
		connections:          mat.NewDense(config.ColumnCount, config.InputWidth, nil),
		permanences:          mat.NewDense(config.ColumnCount, config.InputWidth, nil),
		connectedSynapses:    mat.NewDense(config.ColumnCount, config.InputWidth, nil),
		activeDutyCycles:     make([]float64, config.ColumnCount),
		overlapDutyCycles:    make([]float64, config.ColumnCount),
		minOverlapThresholds: make([]int, config.ColumnCount),
		boostFactors:         make([]float64, config.ColumnCount),
		iterationNum:         0,
		lastUpdateTime:       time.Now(),
		metrics:              htm.NewSpatialPoolerMetrics(),
	}

	// Initialize random number generator
	if config.IsDeterministic() {
		sp.rng = rand.New(rand.NewSource(42)) // Fixed seed for deterministic behavior
	} else {
		sp.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	// Initialize spatial pooler state
	if err := sp.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize spatial pooler: %w", err)
	}

	return sp, nil
}

// Process transforms encoder output into normalized SDR using spatial pooling
func (sp *SpatialPooler) Process(input *htm.PoolingInput) (*htm.PoolingResult, error) {
	startTime := time.Now()

	// Validate input
	if err := input.Validate(); err != nil {
		sp.metrics.RecordError(htm.PoolingErrorInvalidInput)
		return nil, err
	}

	// Convert encoder output to input vector
	inputVector := sp.createInputVector(input.EncoderOutput)

	// Phase 1: Calculate overlap scores
	overlapScores, err := sp.calculateOverlap(inputVector)
	if err != nil {
		sp.metrics.RecordError(htm.PoolingErrorProcessing)
		return nil, fmt.Errorf("overlap calculation failed: %w", err)
	}

	// Phase 2: Apply boosting (if learning enabled)
	if input.LearningEnabled && sp.config.IsLearningEnabled() {
		sp.applyBoosting(overlapScores)
	}

	// Phase 3: Competitive inhibition to select winners
	activeColumns, err := sp.competitiveInhibition(overlapScores)
	if err != nil {
		sp.metrics.RecordError(htm.PoolingErrorProcessing)
		return nil, fmt.Errorf("competitive inhibition failed: %w", err)
	}

	// Phase 4: Learning (update permanences if enabled)
	learningOccurred := false
	if input.LearningEnabled && sp.config.IsLearningEnabled() {
		if err := sp.adaptSynapses(inputVector, activeColumns); err != nil {
			sp.metrics.RecordError(htm.PoolingErrorLearning)
			return nil, fmt.Errorf("learning failed: %w", err)
		}
		learningOccurred = true
		sp.updateDutyCycles(activeColumns, overlapScores)
	}

	// Phase 5: Create output SDR
	outputSDR, err := sp.createOutputSDR(activeColumns)
	if err != nil {
		sp.metrics.RecordError(htm.PoolingErrorProcessing)
		return nil, fmt.Errorf("output SDR creation failed: %w", err)
	}

	// Calculate metrics
	processingTime := time.Since(startTime).Milliseconds()
	avgOverlap := sp.calculateAverageOverlap(overlapScores, activeColumns)
	boostingApplied := input.LearningEnabled && sp.config.BoostStrength > 0

	// Create result
	result := &htm.PoolingResult{
		NormalizedSDR:    *outputSDR,
		InputID:          input.InputID,
		ProcessingTime:   processingTime,
		ActiveColumns:    activeColumns,
		AvgOverlap:       avgOverlap,
		SparsityLevel:    outputSDR.Sparsity,
		LearningOccurred: learningOccurred,
		BoostingApplied:  boostingApplied,
	}

	// Validate result
	if err := result.Validate(); err != nil {
		sp.metrics.RecordError(htm.PoolingErrorProcessing)
		return nil, err
	}

	// Record metrics
	sp.metrics.RecordProcessing(processingTime, outputSDR.Sparsity, learningOccurred, boostingApplied)
	sp.iterationNum++

	return result, nil
}

// GetConfiguration returns the current spatial pooler configuration
func (sp *SpatialPooler) GetConfiguration() *htm.SpatialPoolerConfig {
	// Return a copy to prevent external modification
	configCopy := *sp.config
	return &configCopy
}

// GetMetrics returns current performance and behavioral metrics
func (sp *SpatialPooler) GetMetrics() *htm.SpatialPoolerMetrics {
	// Return a copy to prevent external modification
	metricsCopy := *sp.metrics
	return &metricsCopy
}

// UpdateConfiguration updates the spatial pooler configuration
func (sp *SpatialPooler) UpdateConfiguration(newConfig *htm.SpatialPoolerConfig) error {
	if err := newConfig.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Validate that core dimensions haven't changed
	if newConfig.InputWidth != sp.config.InputWidth {
		return errors.New("cannot change input width after initialization")
	}
	if newConfig.ColumnCount != sp.config.ColumnCount {
		return errors.New("cannot change column count after initialization")
	}

	sp.config = newConfig
	return nil
}

// Private methods

// initialize sets up the initial state of the spatial pooler
func (sp *SpatialPooler) initialize() error {
	// Initialize permanences with random values around 0.5
	rows, cols := sp.permanences.Dims()
	permanenceData := make([]float64, rows*cols)

	for i := range permanenceData {
		// Random permanence values centered around 0.5 with small variation
		permanenceData[i] = 0.5 + (sp.rng.Float64()-0.5)*0.2
	}
	sp.permanences = mat.NewDense(rows, cols, permanenceData)

	// Update connected synapses based on initial permanences
	sp.updateConnectedSynapses()

	// Initialize boost factors to 1.0
	for i := range sp.boostFactors {
		sp.boostFactors[i] = 1.0
	}

	// Initialize minimum overlap thresholds
	for i := range sp.minOverlapThresholds {
		sp.minOverlapThresholds[i] = sp.config.MinOverlapThreshold
	}

	return nil
}

// createInputVector converts encoder output to dense input vector
func (sp *SpatialPooler) createInputVector(encoderOutput htm.EncoderOutput) []float64 {
	inputVector := make([]float64, encoderOutput.Width)
	for _, bit := range encoderOutput.ActiveBits {
		inputVector[bit] = 1.0
	}
	return inputVector
}

// calculateOverlap computes overlap scores for all columns
func (sp *SpatialPooler) calculateOverlap(inputVector []float64) ([]float64, error) {
	rows, cols := sp.connectedSynapses.Dims()
	if cols != len(inputVector) {
		return nil, fmt.Errorf("input vector size (%d) doesn't match expected size (%d)", len(inputVector), cols)
	}

	overlapScores := make([]float64, rows)

	// For each column, calculate overlap with input
	for col := 0; col < rows; col++ {
		overlap := 0.0
		for input := 0; input < cols; input++ {
			if sp.connectedSynapses.At(col, input) > 0 && inputVector[input] > 0 {
				overlap++
			}
		}

		// Apply minimum overlap threshold
		if overlap < float64(sp.minOverlapThresholds[col]) {
			overlap = 0.0
		}

		overlapScores[col] = overlap
	}

	return overlapScores, nil
}

// applyBoosting applies boosting factors to overlap scores
func (sp *SpatialPooler) applyBoosting(overlapScores []float64) {
	for i, overlap := range overlapScores {
		if overlap > 0 {
			overlapScores[i] = overlap * sp.boostFactors[i]
		}
	}
}

// competitiveInhibition selects winning columns based on overlap scores
func (sp *SpatialPooler) competitiveInhibition(overlapScores []float64) ([]int, error) {
	targetActiveColumns := sp.config.GetExpectedActiveColumns()

	// Create column-score pairs for sorting
	type columnScore struct {
		column int
		score  float64
	}

	scores := make([]columnScore, len(overlapScores))
	for i, score := range overlapScores {
		scores[i] = columnScore{column: i, score: score}
	}

	// Sort by score (descending)
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].score == scores[j].score {
			// Tie-breaking: deterministic vs randomized
			if sp.config.IsDeterministic() {
				return scores[i].column < scores[j].column // Consistent tie-breaking
			} else {
				return sp.rng.Float64() < 0.5 // Random tie-breaking
			}
		}
		return scores[i].score > scores[j].score
	})

	// Select top columns with non-zero scores
	activeColumns := make([]int, 0, targetActiveColumns)
	for i := 0; i < len(scores) && len(activeColumns) < targetActiveColumns; i++ {
		if scores[i].score > 0 {
			activeColumns = append(activeColumns, scores[i].column)
		}
	}

	// Ensure minimum number of active columns
	if len(activeColumns) == 0 && targetActiveColumns > 0 {
		// Emergency fallback: activate random columns
		for len(activeColumns) < targetActiveColumns && len(activeColumns) < len(scores) {
			col := scores[len(activeColumns)].column
			activeColumns = append(activeColumns, col)
		}
	}

	// Sort active columns for consistent output
	sort.Ints(activeColumns)

	return activeColumns, nil
}

// adaptSynapses updates permanence values based on learning rules
func (sp *SpatialPooler) adaptSynapses(inputVector []float64, activeColumns []int) error {
	permanenceIncrement := sp.config.LearningRate * 0.1  // 10% of learning rate for increments
	permanenceDecrement := sp.config.LearningRate * 0.05 // 5% of learning rate for decrements

	for _, col := range activeColumns {
		for input := 0; input < len(inputVector); input++ {
			currentPermanence := sp.permanences.At(col, input)
			var newPermanence float64

			if inputVector[input] > 0 {
				// Strengthen connection to active input
				newPermanence = math.Min(1.0, currentPermanence+permanenceIncrement)
			} else {
				// Weaken connection to inactive input
				newPermanence = math.Max(0.0, currentPermanence-permanenceDecrement)
			}

			sp.permanences.Set(col, input, newPermanence)
		}
	}

	// Update connected synapses based on new permanences
	sp.updateConnectedSynapses()

	return nil
}

// updateConnectedSynapses updates the binary connected synapses matrix
func (sp *SpatialPooler) updateConnectedSynapses() {
	rows, cols := sp.permanences.Dims()
	connectedThreshold := 0.5 // Standard HTM permanence threshold

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if sp.permanences.At(row, col) >= connectedThreshold {
				sp.connectedSynapses.Set(row, col, 1.0)
			} else {
				sp.connectedSynapses.Set(row, col, 0.0)
			}
		}
	}
}

// updateDutyCycles updates duty cycles for boosting
func (sp *SpatialPooler) updateDutyCycles(activeColumns []int, overlapScores []float64) {
	dutyCycleDecay := 0.01 // Small decay factor for moving average

	// Update active duty cycles
	for i := range sp.activeDutyCycles {
		isActive := false
		for _, col := range activeColumns {
			if col == i {
				isActive = true
				break
			}
		}

		if isActive {
			sp.activeDutyCycles[i] = sp.activeDutyCycles[i]*(1-dutyCycleDecay) + dutyCycleDecay
		} else {
			sp.activeDutyCycles[i] = sp.activeDutyCycles[i] * (1 - dutyCycleDecay)
		}
	}

	// Update overlap duty cycles
	for i, overlap := range overlapScores {
		hasOverlap := overlap > 0
		if hasOverlap {
			sp.overlapDutyCycles[i] = sp.overlapDutyCycles[i]*(1-dutyCycleDecay) + dutyCycleDecay
		} else {
			sp.overlapDutyCycles[i] = sp.overlapDutyCycles[i] * (1 - dutyCycleDecay)
		}
	}

	// Update boost factors based on duty cycles
	sp.updateBoostFactors()
}

// updateBoostFactors updates boosting factors based on duty cycles
func (sp *SpatialPooler) updateBoostFactors() {
	if sp.config.BoostStrength == 0 {
		return // No boosting
	}

	targetDensity := sp.config.LocalAreaDensity

	for i := range sp.boostFactors {
		dutyCycle := sp.activeDutyCycles[i]

		// Calculate boost factor
		if dutyCycle > 0 {
			boostFactor := math.Exp((targetDensity - dutyCycle) * sp.config.BoostStrength)
			sp.boostFactors[i] = math.Min(sp.config.MaxBoost, math.Max(1.0, boostFactor))
		} else {
			sp.boostFactors[i] = sp.config.MaxBoost
		}
	}
}

// createOutputSDR creates the output SDR from active columns
func (sp *SpatialPooler) createOutputSDR(activeColumns []int) (*sdr.SDR, error) {
	// Map columns to SDR bits (1:1 mapping for spatial pooler)
	return sdr.NewSDR(sp.config.ColumnCount, activeColumns)
}

// calculateAverageOverlap calculates average overlap for active columns
func (sp *SpatialPooler) calculateAverageOverlap(overlapScores []float64, activeColumns []int) float64 {
	if len(activeColumns) == 0 {
		return 0.0
	}

	total := 0.0
	for _, col := range activeColumns {
		total += overlapScores[col]
	}

	return total / float64(len(activeColumns))
}
