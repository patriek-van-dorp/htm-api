package spatial

import (
	"fmt"
	"math"

	"github.com/htm-project/neural-api/internal/domain/htm"
)

// ParameterManager handles spatial pooler parameter management and tuning
type ParameterManager struct {
	config *htm.SpatialPoolerConfig
}

// NewParameterManager creates a new parameter manager
func NewParameterManager(config *htm.SpatialPoolerConfig) *ParameterManager {
	return &ParameterManager{config: config}
}

// ValidateParameterConsistency checks for parameter consistency issues
func (pm *ParameterManager) ValidateParameterConsistency() []string {
	issues := []string{}

	// Check sparsity consistency
	expectedActive := float64(pm.config.ColumnCount) * pm.config.SparsityRatio
	if expectedActive < 1 {
		issues = append(issues, "sparsity ratio too low: would produce < 1 active column")
	}

	// Check learning rate consistency
	if pm.config.LearningEnabled && pm.config.LearningRate == 0 {
		issues = append(issues, "learning enabled but learning rate is zero")
	}

	// Check boosting consistency
	if pm.config.BoostStrength > 0 && pm.config.MaxBoost <= 1 {
		issues = append(issues, "boost strength > 0 but max boost <= 1")
	}

	// Check inhibition radius vs column count
	if pm.config.InhibitionRadius >= pm.config.ColumnCount {
		issues = append(issues, "inhibition radius >= column count")
	}

	// Check semantic thresholds
	if pm.config.SemanticThresholds.SimilarInputMinOverlap <= pm.config.SemanticThresholds.DifferentInputMaxOverlap {
		issues = append(issues, "similar input threshold <= different input threshold")
	}

	return issues
}

// OptimizeForThroughput adjusts parameters for maximum throughput
func (pm *ParameterManager) OptimizeForThroughput() *htm.SpatialPoolerConfig {
	optimized := *pm.config

	// Disable learning for maximum speed
	optimized.LearningEnabled = false

	// Use deterministic mode for consistency
	optimized.Mode = htm.SpatialPoolerModeDeterministic

	// Minimize boosting overhead
	optimized.BoostStrength = 0.0

	// Set conservative processing time limit
	optimized.MaxProcessingTimeMs = 5

	return &optimized
}

// OptimizeForAccuracy adjusts parameters for maximum accuracy
func (pm *ParameterManager) OptimizeForAccuracy() *htm.SpatialPoolerConfig {
	optimized := *pm.config

	// Enable learning
	optimized.LearningEnabled = true

	// Use moderate learning rate
	optimized.LearningRate = 0.1

	// Enable boosting for better coverage
	optimized.BoostStrength = 0.1
	optimized.MaxBoost = 2.0

	// Tighter semantic thresholds
	optimized.SemanticThresholds.SimilarInputMinOverlap = 0.6
	optimized.SemanticThresholds.DifferentInputMaxOverlap = 0.1

	return &optimized
}

// CalculateOptimalColumnCount calculates optimal column count for given input width
func (pm *ParameterManager) CalculateOptimalColumnCount(inputWidth int, targetSparsity float64) int {
	// HTM rule of thumb: column count should be 1.5-2x input width for good coverage
	// while maintaining sparsity constraints

	minColumns := int(float64(inputWidth) * 1.5)
	maxColumns := int(float64(inputWidth) * 2.0)

	// Ensure we can achieve target sparsity with reasonable active count
	minActiveColumns := 20  // Minimum for good representation
	maxActiveColumns := 200 // Maximum for computational efficiency

	optimalColumns := minColumns
	for cols := minColumns; cols <= maxColumns; cols += 10 {
		activeCount := int(float64(cols) * targetSparsity)
		if activeCount >= minActiveColumns && activeCount <= maxActiveColumns {
			optimalColumns = cols
			break
		}
	}

	return optimalColumns
}

// EstimateMemoryUsage estimates memory usage for given parameters
func (pm *ParameterManager) EstimateMemoryUsage() *MemoryEstimate {
	cols := pm.config.ColumnCount
	inputs := pm.config.InputWidth

	// Matrix storage (float64 = 8 bytes)
	permanencesBytes := cols * inputs * 8
	connectionsBytes := cols * inputs * 8
	connectedSynapsesBytes := cols * inputs * 8

	// Column state arrays
	dutyCyclesBytes := cols * 8 * 2 // active + overlap duty cycles
	boostFactorsBytes := cols * 8
	thresholdsBytes := cols * 4 // int32

	totalBytes := permanencesBytes + connectionsBytes + connectedSynapsesBytes +
		dutyCyclesBytes + boostFactorsBytes + thresholdsBytes

	return &MemoryEstimate{
		TotalBytes:       totalBytes,
		TotalMB:          float64(totalBytes) / (1024 * 1024),
		PermanencesBytes: permanencesBytes,
		StateBytes:       dutyCyclesBytes + boostFactorsBytes + thresholdsBytes,
		Breakdown: map[string]int{
			"permanences":        permanencesBytes,
			"connections":        connectionsBytes,
			"connected_synapses": connectedSynapsesBytes,
			"duty_cycles":        dutyCyclesBytes,
			"boost_factors":      boostFactorsBytes,
			"thresholds":         thresholdsBytes,
		},
	}
}

// MemoryEstimate contains memory usage estimation
type MemoryEstimate struct {
	TotalBytes       int            `json:"total_bytes"`
	TotalMB          float64        `json:"total_mb"`
	PermanencesBytes int            `json:"permanences_bytes"`
	StateBytes       int            `json:"state_bytes"`
	Breakdown        map[string]int `json:"breakdown"`
}

// CalculateProcessingComplexity estimates computational complexity
func (pm *ParameterManager) CalculateProcessingComplexity() *ComplexityEstimate {
	cols := pm.config.ColumnCount
	inputs := pm.config.InputWidth

	// O(cols * inputs) for overlap calculation
	overlapOps := cols * inputs

	// O(cols * log(cols)) for sorting in competitive inhibition
	inhibitionOps := int(float64(cols) * math.Log2(float64(cols)))

	// O(activeColumns * inputs) for learning (if enabled)
	activeColumns := pm.config.GetExpectedActiveColumns()
	learningOps := 0
	if pm.config.LearningEnabled {
		learningOps = activeColumns * inputs
	}

	totalOps := overlapOps + inhibitionOps + learningOps

	return &ComplexityEstimate{
		TotalOperations:       totalOps,
		OverlapOps:            overlapOps,
		InhibitionOps:         inhibitionOps,
		LearningOps:           learningOps,
		EstimatedMicroseconds: int64(float64(totalOps) / 1000), // Rough estimate
	}
}

// ComplexityEstimate contains computational complexity analysis
type ComplexityEstimate struct {
	TotalOperations       int   `json:"total_operations"`
	OverlapOps            int   `json:"overlap_ops"`
	InhibitionOps         int   `json:"inhibition_ops"`
	LearningOps           int   `json:"learning_ops"`
	EstimatedMicroseconds int64 `json:"estimated_microseconds"`
}

// AutoTuneParameters automatically adjusts parameters based on performance feedback
func (pm *ParameterManager) AutoTuneParameters(metrics *htm.SpatialPoolerMetrics, targetPerformance *PerformanceTarget) *htm.SpatialPoolerConfig {
	tuned := *pm.config

	// Adjust based on processing time
	if metrics.AverageProcessingTime > targetPerformance.MaxProcessingTimeMs {
		// Reduce complexity
		if tuned.LearningEnabled && targetPerformance.AccuracyPriority < 0.8 {
			tuned.LearningEnabled = false
		}
		tuned.BoostStrength *= 0.9 // Reduce boosting overhead
	}

	// Adjust based on sparsity
	if metrics.AverageSparsity < targetPerformance.MinSparsity {
		tuned.SparsityRatio = math.Min(0.05, tuned.SparsityRatio*1.1)
	} else if metrics.AverageSparsity > targetPerformance.MaxSparsity {
		tuned.SparsityRatio = math.Max(0.02, tuned.SparsityRatio*0.9)
	}

	// Adjust learning rate based on learning effectiveness
	if tuned.LearningEnabled && metrics.LearningIterations > 0 {
		if metrics.BoostingEvents > metrics.LearningIterations/10 {
			// Too much boosting, increase learning rate
			tuned.LearningRate = math.Min(1.0, tuned.LearningRate*1.1)
		}
	}

	return &tuned
}

// PerformanceTarget defines performance optimization targets
type PerformanceTarget struct {
	MaxProcessingTimeMs int64   `json:"max_processing_time_ms"`
	MinSparsity         float64 `json:"min_sparsity"`
	MaxSparsity         float64 `json:"max_sparsity"`
	AccuracyPriority    float64 `json:"accuracy_priority"` // 0.0 (speed) to 1.0 (accuracy)
}

// GetParameterRecommendations provides parameter recommendations for different use cases
func (pm *ParameterManager) GetParameterRecommendations(useCase string) (*htm.SpatialPoolerConfig, error) {
	base := htm.DefaultSpatialPoolerConfig()

	switch useCase {
	case "high_throughput":
		base.LearningEnabled = false
		base.Mode = htm.SpatialPoolerModeDeterministic
		base.BoostStrength = 0.0
		base.MaxProcessingTimeMs = 5

	case "high_accuracy":
		base.LearningEnabled = true
		base.LearningRate = 0.15
		base.BoostStrength = 0.2
		base.MaxBoost = 3.0
		base.SemanticThresholds.SimilarInputMinOverlap = 0.7
		base.SemanticThresholds.DifferentInputMaxOverlap = 0.05

	case "balanced":
		base.LearningEnabled = true
		base.LearningRate = 0.1
		base.BoostStrength = 0.1
		base.MaxBoost = 2.0

	case "memory_efficient":
		base.ColumnCount = base.InputWidth // 1:1 ratio
		base.SparsityRatio = 0.02          // Minimum sparsity
		base.BoostStrength = 0.0           // No boosting overhead

	default:
		return nil, fmt.Errorf("unknown use case: %s", useCase)
	}

	return base, nil
}
