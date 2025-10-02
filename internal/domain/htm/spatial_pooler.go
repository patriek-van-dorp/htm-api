package htm

import (
	"fmt"
	"strings"

	"github.com/htm-project/neural-api/internal/cortical/sdr"
)

// SpatialPoolerMode represents the processing mode for the spatial pooler
type SpatialPoolerMode string

const (
	// SpatialPoolerModeDeterministic - Identical inputs always produce identical outputs
	SpatialPoolerModeDeterministic SpatialPoolerMode = "deterministic"
	// SpatialPoolerModeRandomized - Controlled randomness for learning purposes
	SpatialPoolerModeRandomized SpatialPoolerMode = "randomized"
)

// IsValid checks if the spatial pooler mode is valid
func (m SpatialPoolerMode) IsValid() bool {
	switch m {
	case SpatialPoolerModeDeterministic, SpatialPoolerModeRandomized:
		return true
	default:
		return false
	}
}

// String returns the string representation of the mode
func (m SpatialPoolerMode) String() string {
	return string(m)
}

// ParseSpatialPoolerMode parses a string into a SpatialPoolerMode
func ParseSpatialPoolerMode(s string) (SpatialPoolerMode, error) {
	mode := SpatialPoolerMode(strings.ToLower(strings.TrimSpace(s)))
	if !mode.IsValid() {
		return "", fmt.Errorf("invalid spatial pooler mode: %s", s)
	}
	return mode, nil
}

// PoolingError represents errors that can occur during spatial pooling operations
type PoolingError struct {
	ErrorType   PoolingErrorType `json:"error_type"`
	Message     string           `json:"message"`
	InputID     string           `json:"input_id,omitempty"`
	ConfigField string           `json:"config_field,omitempty"`
}

// Error implements the error interface
func (e *PoolingError) Error() string {
	if e.InputID != "" {
		return fmt.Sprintf("%s: %s (input: %s)", e.ErrorType, e.Message, e.InputID)
	}
	if e.ConfigField != "" {
		return fmt.Sprintf("%s: %s (field: %s)", e.ErrorType, e.Message, e.ConfigField)
	}
	return fmt.Sprintf("%s: %s", e.ErrorType, e.Message)
}

// PoolingErrorType represents the category of pooling error
type PoolingErrorType string

const (
	// PoolingErrorInvalidInput - Input validation failed
	PoolingErrorInvalidInput PoolingErrorType = "invalid_input"
	// PoolingErrorConfiguration - Invalid spatial pooler configuration
	PoolingErrorConfiguration PoolingErrorType = "configuration_error"
	// PoolingErrorProcessing - Error during spatial pooling computation
	PoolingErrorProcessing PoolingErrorType = "processing_error"
	// PoolingErrorPerformance - Processing exceeded time/memory constraints
	PoolingErrorPerformance PoolingErrorType = "performance_error"
	// PoolingErrorLearning - Error during learning rule application
	PoolingErrorLearning PoolingErrorType = "learning_error"
)

// IsValid checks if the pooling error type is valid
func (e PoolingErrorType) IsValid() bool {
	switch e {
	case PoolingErrorInvalidInput, PoolingErrorConfiguration, PoolingErrorProcessing,
		PoolingErrorPerformance, PoolingErrorLearning:
		return true
	default:
		return false
	}
}

// String returns the string representation of the error type
func (e PoolingErrorType) String() string {
	return string(e)
}

// NewPoolingError creates a new pooling error
func NewPoolingError(errorType PoolingErrorType, message string) *PoolingError {
	return &PoolingError{
		ErrorType: errorType,
		Message:   message,
	}
}

// NewPoolingErrorWithInput creates a new pooling error with input ID
func NewPoolingErrorWithInput(errorType PoolingErrorType, message, inputID string) *PoolingError {
	return &PoolingError{
		ErrorType: errorType,
		Message:   message,
		InputID:   inputID,
	}
}

// NewPoolingErrorWithField creates a new pooling error with config field
func NewPoolingErrorWithField(errorType PoolingErrorType, message, configField string) *PoolingError {
	return &PoolingError{
		ErrorType:   errorType,
		Message:     message,
		ConfigField: configField,
	}
}

// PoolingInput represents input structure for spatial pooler processing
type PoolingInput struct {
	EncoderOutput   EncoderOutput          `json:"encoder_output" validate:"required"`
	InputWidth      int                    `json:"input_width" validate:"required,gt=0"`
	InputID         string                 `json:"input_id" validate:"required"`
	LearningEnabled bool                   `json:"learning_enabled"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// EncoderOutput represents raw bit array from sensor encoder
type EncoderOutput struct {
	Width      int     `json:"width" validate:"required,gt=0"`
	ActiveBits []int   `json:"active_bits" validate:"required"`
	Sparsity   float64 `json:"sparsity" validate:"gte=0,lte=1"`
}

// Validate validates the pooling input
func (p *PoolingInput) Validate() error {
	if p.EncoderOutput.Width <= 0 {
		return NewPoolingError(PoolingErrorInvalidInput, "encoder output width must be positive")
	}

	if len(p.EncoderOutput.ActiveBits) == 0 {
		return NewPoolingErrorWithInput(PoolingErrorInvalidInput, "encoder output must have active bits", p.InputID)
	}

	if p.InputWidth != p.EncoderOutput.Width {
		return NewPoolingErrorWithInput(PoolingErrorInvalidInput,
			fmt.Sprintf("input width (%d) must match encoder output width (%d)", p.InputWidth, p.EncoderOutput.Width),
			p.InputID)
	}

	if p.InputID == "" {
		return NewPoolingError(PoolingErrorInvalidInput, "input ID cannot be empty")
	}

	// Validate active bits are within valid range
	for _, bit := range p.EncoderOutput.ActiveBits {
		if bit < 0 || bit >= p.EncoderOutput.Width {
			return NewPoolingErrorWithInput(PoolingErrorInvalidInput,
				fmt.Sprintf("active bit %d is out of range [0, %d)", bit, p.EncoderOutput.Width),
				p.InputID)
		}
	}

	return nil
}

// PoolingResult represents output structure containing true SDR produced by spatial pooler
type PoolingResult struct {
	NormalizedSDR    sdr.SDR `json:"normalized_sdr"`
	InputID          string  `json:"input_id"`
	ProcessingTime   int64   `json:"processing_time_ms"` // milliseconds
	ActiveColumns    []int   `json:"active_columns"`
	AvgOverlap       float64 `json:"avg_overlap"`
	SparsityLevel    float64 `json:"sparsity_level"`
	LearningOccurred bool    `json:"learning_occurred"`
	BoostingApplied  bool    `json:"boosting_applied"`
}

// Validate validates the pooling result
func (r *PoolingResult) Validate() error {
	// Validate sparsity is in HTM range (2-5%)
	if r.SparsityLevel < 0.02 || r.SparsityLevel > 0.05 {
		return NewPoolingError(PoolingErrorProcessing,
			fmt.Sprintf("sparsity level %.4f is outside HTM range [0.02, 0.05]", r.SparsityLevel))
	}

	// Validate processing time meets performance requirements (<10ms)
	if r.ProcessingTime > 10 {
		return NewPoolingError(PoolingErrorPerformance,
			fmt.Sprintf("processing time %dms exceeds 10ms requirement", r.ProcessingTime))
	}

	// Validate active columns are sorted and within valid range
	if len(r.ActiveColumns) == 0 {
		return NewPoolingError(PoolingErrorProcessing, "result must have active columns")
	}

	for i, col := range r.ActiveColumns {
		if col < 0 {
			return NewPoolingError(PoolingErrorProcessing,
				fmt.Sprintf("active column %d is negative", col))
		}
		if i > 0 && col <= r.ActiveColumns[i-1] {
			return NewPoolingError(PoolingErrorProcessing, "active columns must be sorted")
		}
	}

	// Validate sparsity matches SDR
	expectedSparsity := float64(len(r.NormalizedSDR.ActiveBits)) / float64(r.NormalizedSDR.Width)
	if abs(r.SparsityLevel-expectedSparsity) > 0.001 {
		return NewPoolingError(PoolingErrorProcessing,
			fmt.Sprintf("sparsity level %.4f does not match SDR sparsity %.4f", r.SparsityLevel, expectedSparsity))
	}

	return nil
}

// SpatialPoolerMetrics represents performance and behavioral metrics
type SpatialPoolerMetrics struct {
	TotalProcessed           int64                      `json:"total_processed"`
	AverageProcessingTime    int64                      `json:"average_processing_time_ms"` // milliseconds
	LearningIterations       int64                      `json:"learning_iterations"`
	ColumnUsageDistribution  []float64                  `json:"column_usage_distribution,omitempty"`
	AverageSparsity          float64                    `json:"average_sparsity"`
	OverlapScoreDistribution []float64                  `json:"overlap_score_distribution,omitempty"`
	BoostingEvents           int64                      `json:"boosting_events"`
	ErrorCounts              map[PoolingErrorType]int64 `json:"error_counts"`
}

// NewSpatialPoolerMetrics creates a new metrics instance
func NewSpatialPoolerMetrics() *SpatialPoolerMetrics {
	return &SpatialPoolerMetrics{
		ErrorCounts: make(map[PoolingErrorType]int64),
	}
}

// RecordProcessing records a successful processing operation
func (m *SpatialPoolerMetrics) RecordProcessing(processingTime int64, sparsity float64, learningOccurred bool, boostingApplied bool) {
	m.TotalProcessed++

	// Update average processing time
	if m.TotalProcessed == 1 {
		m.AverageProcessingTime = processingTime
	} else {
		// Running average: new_avg = old_avg + (new_value - old_avg) / count
		m.AverageProcessingTime = m.AverageProcessingTime + (processingTime-m.AverageProcessingTime)/m.TotalProcessed
	}

	// Update average sparsity
	if m.TotalProcessed == 1 {
		m.AverageSparsity = sparsity
	} else {
		m.AverageSparsity = m.AverageSparsity + (sparsity-m.AverageSparsity)/float64(m.TotalProcessed)
	}

	if learningOccurred {
		m.LearningIterations++
	}

	if boostingApplied {
		m.BoostingEvents++
	}
}

// RecordError records an error occurrence
func (m *SpatialPoolerMetrics) RecordError(errorType PoolingErrorType) {
	if m.ErrorCounts == nil {
		m.ErrorCounts = make(map[PoolingErrorType]int64)
	}
	m.ErrorCounts[errorType]++
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// SpatialPoolerConfig represents configuration parameters for spatial pooler
type SpatialPoolerConfig struct {
	// Core dimensions
	InputWidth  int `json:"input_width" validate:"required,gt=0"`
	ColumnCount int `json:"column_count" validate:"required,gt=0"`

	// Sparsity control (HTM requirement: 2-5%)
	SparsityRatio float64 `json:"sparsity_ratio" validate:"required,gte=0.02,lte=0.05"`

	// Processing mode
	Mode SpatialPoolerMode `json:"mode" validate:"required"`

	// Learning parameters
	LearningEnabled bool    `json:"learning_enabled"`
	LearningRate    float64 `json:"learning_rate" validate:"gte=0,lte=1"`
	MaxBoost        float64 `json:"max_boost" validate:"gte=1,lte=10"`
	BoostStrength   float64 `json:"boost_strength" validate:"gte=0,lte=1"`

	// Competitive inhibition
	InhibitionRadius    int     `json:"inhibition_radius" validate:"gte=0"`
	LocalAreaDensity    float64 `json:"local_area_density" validate:"gte=0,lte=1"`
	MinOverlapThreshold int     `json:"min_overlap_threshold" validate:"gte=0"`

	// Performance constraints
	MaxProcessingTimeMs int `json:"max_processing_time_ms" validate:"gte=1,lte=10"` // <10ms requirement

	// Semantic similarity preservation
	SemanticThresholds SemanticThresholds `json:"semantic_thresholds"`
}

// SemanticThresholds represents thresholds for semantic similarity preservation
type SemanticThresholds struct {
	SimilarInputMinOverlap   float64 `json:"similar_input_min_overlap" validate:"gte=0.3,lte=0.7"` // 30-70% for similar inputs
	DifferentInputMaxOverlap float64 `json:"different_input_max_overlap" validate:"gte=0,lte=0.2"` // <20% for different inputs
}

// DefaultSpatialPoolerConfig returns a default configuration following HTM principles
func DefaultSpatialPoolerConfig() *SpatialPoolerConfig {
	return &SpatialPoolerConfig{
		InputWidth:          1024,
		ColumnCount:         2048,
		SparsityRatio:       0.02, // 2% - HTM standard
		Mode:                SpatialPoolerModeDeterministic,
		LearningEnabled:     true,
		LearningRate:        0.1,
		MaxBoost:            2.0,
		BoostStrength:       0.0,
		InhibitionRadius:    16,
		LocalAreaDensity:    0.02, // Matches sparsity ratio
		MinOverlapThreshold: 1,
		MaxProcessingTimeMs: 10, // Performance requirement
		SemanticThresholds: SemanticThresholds{
			SimilarInputMinOverlap:   0.5, // 50% overlap for similar inputs
			DifferentInputMaxOverlap: 0.1, // 10% overlap for different inputs
		},
	}
}

// Validate validates the spatial pooler configuration
func (c *SpatialPoolerConfig) Validate() error {
	// Basic validation
	if c.InputWidth <= 0 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration, "input width must be positive", "input_width")
	}

	if c.ColumnCount <= 0 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration, "column count must be positive", "column_count")
	}

	// HTM sparsity requirement
	if c.SparsityRatio < 0.02 || c.SparsityRatio > 0.05 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration,
			fmt.Sprintf("sparsity ratio %.4f must be between 0.02 and 0.05 (HTM requirement)", c.SparsityRatio),
			"sparsity_ratio")
	}

	// Mode validation
	if !c.Mode.IsValid() {
		return NewPoolingErrorWithField(PoolingErrorConfiguration,
			fmt.Sprintf("invalid mode: %s", c.Mode), "mode")
	}

	// Learning parameters
	if c.LearningRate < 0 || c.LearningRate > 1 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration, "learning rate must be between 0 and 1", "learning_rate")
	}

	if c.MaxBoost < 1 || c.MaxBoost > 10 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration, "max boost must be between 1 and 10", "max_boost")
	}

	if c.BoostStrength < 0 || c.BoostStrength > 1 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration, "boost strength must be between 0 and 1", "boost_strength")
	}

	// Inhibition parameters
	if c.InhibitionRadius < 0 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration, "inhibition radius cannot be negative", "inhibition_radius")
	}

	if c.LocalAreaDensity < 0 || c.LocalAreaDensity > 1 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration, "local area density must be between 0 and 1", "local_area_density")
	}

	if c.MinOverlapThreshold < 0 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration, "min overlap threshold cannot be negative", "min_overlap_threshold")
	}

	// Performance constraint
	if c.MaxProcessingTimeMs < 1 || c.MaxProcessingTimeMs > 10 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration, "max processing time must be between 1 and 10 ms", "max_processing_time_ms")
	}

	// Semantic thresholds
	if c.SemanticThresholds.SimilarInputMinOverlap < 0.3 || c.SemanticThresholds.SimilarInputMinOverlap > 0.7 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration,
			"similar input min overlap must be between 0.3 and 0.7", "semantic_thresholds.similar_input_min_overlap")
	}

	if c.SemanticThresholds.DifferentInputMaxOverlap < 0 || c.SemanticThresholds.DifferentInputMaxOverlap > 0.2 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration,
			"different input max overlap must be between 0 and 0.2", "semantic_thresholds.different_input_max_overlap")
	}

	// Cross-validation: ensure thresholds don't overlap
	if c.SemanticThresholds.SimilarInputMinOverlap <= c.SemanticThresholds.DifferentInputMaxOverlap {
		return NewPoolingErrorWithField(PoolingErrorConfiguration,
			"similar input min overlap must be greater than different input max overlap", "semantic_thresholds")
	}

	// Verify sparsity configuration produces reasonable number of active columns
	expectedActiveColumns := int(float64(c.ColumnCount) * c.SparsityRatio)
	if expectedActiveColumns < 1 {
		return NewPoolingErrorWithField(PoolingErrorConfiguration,
			fmt.Sprintf("sparsity ratio %.4f with %d columns produces <1 active column", c.SparsityRatio, c.ColumnCount),
			"sparsity_ratio")
	}

	return nil
}

// GetExpectedActiveColumns calculates the expected number of active columns
func (c *SpatialPoolerConfig) GetExpectedActiveColumns() int {
	return int(float64(c.ColumnCount) * c.SparsityRatio)
}

// GetExpectedSparsity returns the expected sparsity level
func (c *SpatialPoolerConfig) GetExpectedSparsity() float64 {
	return c.SparsityRatio
}

// IsDeterministic returns true if the pooler is configured for deterministic operation
func (c *SpatialPoolerConfig) IsDeterministic() bool {
	return c.Mode == SpatialPoolerModeDeterministic
}

// IsLearningEnabled returns true if learning is enabled
func (c *SpatialPoolerConfig) IsLearningEnabled() bool {
	return c.LearningEnabled
}
