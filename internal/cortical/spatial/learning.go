package spatial

import (
	"fmt"
	"math"
	"sort"

	"github.com/htm-project/neural-api/internal/domain/htm"
	"gonum.org/v1/gonum/mat"
)

// LearningManager implements learning algorithms and synaptic adaptation for spatial pooling
type LearningManager struct {
	config              *htm.SpatialPoolerConfig
	permanenceIncrement float64
	permanenceDecrement float64
	connectedThreshold  float64
	minThreshold        float64
	maxThreshold        float64
}

// NewLearningManager creates a new learning manager
func NewLearningManager(config *htm.SpatialPoolerConfig) *LearningManager {
	return &LearningManager{
		config:              config,
		permanenceIncrement: config.LearningRate * 0.1,  // 10% of learning rate
		permanenceDecrement: config.LearningRate * 0.05, // 5% of learning rate
		connectedThreshold:  0.5,                        // Standard HTM threshold
		minThreshold:        0.0,
		maxThreshold:        1.0,
	}
}

// UpdatePermanences updates synaptic permanences for active columns
func (lm *LearningManager) UpdatePermanences(
	permanences *mat.Dense,
	inputVector []float64,
	activeColumns []int) error {

	if !lm.config.IsLearningEnabled() {
		return nil // Learning disabled
	}

	rows, cols := permanences.Dims()
	if cols != len(inputVector) {
		return fmt.Errorf("input vector size (%d) doesn't match permanence matrix cols (%d)", len(inputVector), cols)
	}

	// Apply Hebbian learning rule for each active column
	for _, col := range activeColumns {
		if col >= rows {
			continue // Skip invalid columns
		}

		for input := 0; input < cols; input++ {
			currentPermanence := permanences.At(col, input)
			var newPermanence float64

			if inputVector[input] > 0 {
				// Strengthen synapses to active inputs (LTP - Long Term Potentiation)
				newPermanence = lm.applySynapticStrengthening(currentPermanence)
			} else {
				// Weaken synapses to inactive inputs (LTD - Long Term Depression)
				newPermanence = lm.applySynapticWeakening(currentPermanence)
			}

			// Clip to valid range
			newPermanence = math.Max(lm.minThreshold, math.Min(lm.maxThreshold, newPermanence))
			permanences.Set(col, input, newPermanence)
		}
	}

	return nil
}

// UpdateBoostFactors updates column boost factors based on duty cycles
func (lm *LearningManager) UpdateBoostFactors(
	activeDutyCycles []float64,
	targetDensity float64) []float64 {

	if lm.config.BoostStrength == 0 {
		// No boosting - return all 1.0
		boostFactors := make([]float64, len(activeDutyCycles))
		for i := range boostFactors {
			boostFactors[i] = 1.0
		}
		return boostFactors
	}

	boostFactors := make([]float64, len(activeDutyCycles))

	for i, dutyCycle := range activeDutyCycles {
		if dutyCycle > 0 {
			// Exponential boosting function
			boostFactor := math.Exp((targetDensity - dutyCycle) * lm.config.BoostStrength)
			boostFactors[i] = math.Min(lm.config.MaxBoost, math.Max(1.0, boostFactor))
		} else {
			// Column never active - apply maximum boost
			boostFactors[i] = lm.config.MaxBoost
		}
	}

	return boostFactors
}

// UpdateDutyCycles updates duty cycles using exponential moving average
func (lm *LearningManager) UpdateDutyCycles(
	currentDutyCycles []float64,
	activeColumns []int,
	decayRate float64) []float64 {

	newDutyCycles := make([]float64, len(currentDutyCycles))

	// Create active column set for fast lookup
	activeSet := make(map[int]bool)
	for _, col := range activeColumns {
		activeSet[col] = true
	}

	// Update each column's duty cycle
	for i := range currentDutyCycles {
		if activeSet[i] {
			// Column was active - increase duty cycle
			newDutyCycles[i] = currentDutyCycles[i]*(1-decayRate) + decayRate
		} else {
			// Column was inactive - decrease duty cycle
			newDutyCycles[i] = currentDutyCycles[i] * (1 - decayRate)
		}
	}

	return newDutyCycles
}

// AdaptMinOverlapThresholds adapts minimum overlap thresholds based on column performance
func (lm *LearningManager) AdaptMinOverlapThresholds(
	currentThresholds []int,
	overlapDutyCycles []float64,
	targetOverlapDensity float64) []int {

	newThresholds := make([]int, len(currentThresholds))

	for i := range currentThresholds {
		overlapDutyCycle := overlapDutyCycles[i]

		if overlapDutyCycle < targetOverlapDensity {
			// Column doesn't overlap enough - decrease threshold
			newThresholds[i] = max(0, currentThresholds[i]-1)
		} else if overlapDutyCycle > targetOverlapDensity*1.5 {
			// Column overlaps too much - increase threshold
			newThresholds[i] = currentThresholds[i] + 1
		} else {
			// Keep current threshold
			newThresholds[i] = currentThresholds[i]
		}
	}

	return newThresholds
}

// HomeostaticPlasticity implements homeostatic plasticity to maintain target sparsity
func (lm *LearningManager) HomeostaticPlasticity(
	permanences *mat.Dense,
	activeDutyCycles []float64,
	targetDensity float64,
	plasticityRate float64) error {

	if plasticityRate == 0 {
		return nil // No homeostatic plasticity
	}

	rows, cols := permanences.Dims()

	for col := 0; col < rows; col++ {
		dutyCycle := activeDutyCycles[col]

		// Calculate scaling factor based on duty cycle deviation
		scalingFactor := 1.0 + plasticityRate*(targetDensity-dutyCycle)

		if scalingFactor != 1.0 {
			// Scale all permanences for this column
			for input := 0; input < cols; input++ {
				currentPermanence := permanences.At(col, input)
				newPermanence := currentPermanence * scalingFactor

				// Clip to valid range
				newPermanence = math.Max(lm.minThreshold, math.Min(lm.maxThreshold, newPermanence))
				permanences.Set(col, input, newPermanence)
			}
		}
	}

	return nil
}

// StructuralPlasticity implements structural plasticity (synapse formation/elimination)
func (lm *LearningManager) StructuralPlasticity(
	permanences *mat.Dense,
	inputVector []float64,
	activeColumns []int,
	newSynapseCount int) error {

	if newSynapseCount <= 0 {
		return nil // No structural plasticity
	}

	rows, cols := permanences.Dims()

	for _, col := range activeColumns {
		if col >= rows {
			continue
		}

		// Find weakest synapses
		synapseStrengths := make([]synapseInfo, cols)
		for input := 0; input < cols; input++ {
			synapseStrengths[input] = synapseInfo{
				index:      input,
				permanence: permanences.At(col, input),
			}
		}

		// Sort by permanence (ascending - weakest first)
		sort.Slice(synapseStrengths, func(i, j int) bool {
			return synapseStrengths[i].permanence < synapseStrengths[j].permanence
		})

		// Replace weakest synapses with new ones to active inputs
		replacedCount := 0
		for _, synapseInfo := range synapseStrengths {
			if replacedCount >= newSynapseCount {
				break
			}

			input := synapseInfo.index
			if inputVector[input] > 0 && synapseInfo.permanence < lm.connectedThreshold {
				// Replace with stronger synapse
				newPermanence := lm.connectedThreshold + 0.1
				permanences.Set(col, input, newPermanence)
				replacedCount++
			}
		}
	}

	return nil
}

// MetaplasticityLearning implements metaplasticity (learning to learn)
func (lm *LearningManager) MetaplasticityLearning(
	permanences *mat.Dense,
	inputVector []float64,
	activeColumns []int,
	priorActivation []float64) error {

	if len(priorActivation) == 0 {
		// No prior activation history - use standard learning
		return lm.UpdatePermanences(permanences, inputVector, activeColumns)
	}

	rows, cols := permanences.Dims()
	if cols != len(inputVector) || rows != len(priorActivation) {
		return fmt.Errorf("dimension mismatch in metaplastic learning")
	}

	// Modulate learning rate based on prior activation
	for _, col := range activeColumns {
		if col >= rows {
			continue
		}

		// Calculate metaplastic modulation factor
		priorActivity := priorActivation[col]
		modulation := lm.calculateMetaplasticModulation(priorActivity)

		for input := 0; input < cols; input++ {
			currentPermanence := permanences.At(col, input)
			var deltaP float64

			if inputVector[input] > 0 {
				deltaP = lm.permanenceIncrement * modulation
			} else {
				deltaP = -lm.permanenceDecrement * modulation
			}

			newPermanence := currentPermanence + deltaP
			newPermanence = math.Max(lm.minThreshold, math.Min(lm.maxThreshold, newPermanence))
			permanences.Set(col, input, newPermanence)
		}
	}

	return nil
}

// LearningMetrics contains metrics about learning performance
type LearningMetrics struct {
	PermanenceChanges  int     `json:"permanence_changes"`
	AveragePermanence  float64 `json:"average_permanence"`
	ConnectedSynapses  int     `json:"connected_synapses"`
	LearningEfficiency float64 `json:"learning_efficiency"`
	StabilityIndex     float64 `json:"stability_index"`
}

// CalculateLearningMetrics calculates learning performance metrics
func (lm *LearningManager) CalculateLearningMetrics(
	permanences *mat.Dense,
	previousPermanences *mat.Dense) *LearningMetrics {

	rows, cols := permanences.Dims()

	permanenceChanges := 0
	totalPermanence := 0.0
	connectedSynapses := 0
	totalChangeMagnitude := 0.0

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			currentPerm := permanences.At(row, col)
			totalPermanence += currentPerm

			if currentPerm >= lm.connectedThreshold {
				connectedSynapses++
			}

			if previousPermanences != nil {
				prevPerm := previousPermanences.At(row, col)
				if currentPerm != prevPerm {
					permanenceChanges++
					totalChangeMagnitude += math.Abs(currentPerm - prevPerm)
				}
			}
		}
	}

	totalSynapses := rows * cols
	averagePermanence := totalPermanence / float64(totalSynapses)

	// Learning efficiency: connected synapses / total changes
	learningEfficiency := 0.0
	if permanenceChanges > 0 {
		learningEfficiency = float64(connectedSynapses) / float64(permanenceChanges)
	}

	// Stability index: 1 - (change magnitude / possible change)
	stabilityIndex := 1.0
	if totalSynapses > 0 {
		maxPossibleChange := float64(totalSynapses) * (lm.maxThreshold - lm.minThreshold)
		if maxPossibleChange > 0 {
			stabilityIndex = 1.0 - (totalChangeMagnitude / maxPossibleChange)
		}
	}

	return &LearningMetrics{
		PermanenceChanges:  permanenceChanges,
		AveragePermanence:  averagePermanence,
		ConnectedSynapses:  connectedSynapses,
		LearningEfficiency: learningEfficiency,
		StabilityIndex:     stabilityIndex,
	}
}

// Private helper methods

func (lm *LearningManager) applySynapticStrengthening(currentPermanence float64) float64 {
	// Non-linear strengthening with diminishing returns near maximum
	increment := lm.permanenceIncrement * (1.0 - currentPermanence)
	return currentPermanence + increment
}

func (lm *LearningManager) applySynapticWeakening(currentPermanence float64) float64 {
	// Linear weakening
	return currentPermanence - lm.permanenceDecrement
}

func (lm *LearningManager) calculateMetaplasticModulation(priorActivity float64) float64 {
	// BCM-like rule: modulation based on deviation from target
	targetActivity := lm.config.LocalAreaDensity
	deviation := priorActivity - targetActivity

	// Sigmoidal modulation function
	return 1.0 / (1.0 + math.Exp(-deviation*5.0))
}

// Helper types and functions

type synapseInfo struct {
	index      int
	permanence float64
}
