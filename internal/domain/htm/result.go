package htm

import "time"

// ProcessingResult represents the output from HTM neural network computation,
// maintaining the same format as input for API chaining.
type ProcessingResult struct {
	ID       string           `json:"id" validate:"required,uuid"`
	Result   [][]float64      `json:"result" validate:"required"`
	Metadata ResultMetadata   `json:"metadata" validate:"required"`
	Status   ProcessingStatus `json:"status" validate:"required"`
}

// ResultMetadata contains processing performance and context information.
type ResultMetadata struct {
	ProcessingTimeMs int64              `json:"processing_time_ms" validate:"min=0"`
	InstanceID       string             `json:"instance_id" validate:"required"`
	AlgorithmVersion string             `json:"algorithm_version" validate:"required"`
	QualityMetrics   map[string]float64 `json:"quality_metrics,omitempty"`
}

// GetProcessingDuration returns the processing time as a time.Duration
func (r *ResultMetadata) GetProcessingDuration() time.Duration {
	return time.Duration(r.ProcessingTimeMs) * time.Millisecond
}

// GetQualityMetric retrieves a quality metric by key
func (r *ResultMetadata) GetQualityMetric(key string) (float64, bool) {
	if r.QualityMetrics == nil {
		return 0.0, false
	}

	value, exists := r.QualityMetrics[key]
	return value, exists
}

// SetQualityMetric sets a quality metric
func (r *ResultMetadata) SetQualityMetric(key string, value float64) {
	if r.QualityMetrics == nil {
		r.QualityMetrics = make(map[string]float64)
	}

	// Ensure quality metrics are between 0.0 and 1.0
	if value < 0.0 {
		value = 0.0
	} else if value > 1.0 {
		value = 1.0
	}

	r.QualityMetrics[key] = value
}

// IsSuccessful returns true if the processing was successful
func (p *ProcessingResult) IsSuccessful() bool {
	return p.Status == StatusSuccess || p.Status == StatusPartialSuccess
}

// HasData returns true if the result contains data
func (p *ProcessingResult) HasData() bool {
	return len(p.Result) > 0
}

// GetDimensions returns the dimensions of the result data
func (p *ProcessingResult) GetDimensions() []int {
	if len(p.Result) == 0 {
		return []int{0, 0}
	}
	return []int{len(p.Result), len(p.Result[0])}
}

// Clone creates a deep copy of the ProcessingResult
func (p *ProcessingResult) Clone() *ProcessingResult {
	// Deep copy the result data
	resultCopy := make([][]float64, len(p.Result))
	for i, row := range p.Result {
		resultCopy[i] = make([]float64, len(row))
		copy(resultCopy[i], row)
	}

	// Deep copy quality metrics
	qualityMetricsCopy := make(map[string]float64)
	for k, v := range p.Metadata.QualityMetrics {
		qualityMetricsCopy[k] = v
	}

	return &ProcessingResult{
		ID:     p.ID,
		Result: resultCopy,
		Metadata: ResultMetadata{
			ProcessingTimeMs: p.Metadata.ProcessingTimeMs,
			InstanceID:       p.Metadata.InstanceID,
			AlgorithmVersion: p.Metadata.AlgorithmVersion,
			QualityMetrics:   qualityMetricsCopy,
		},
		Status: p.Status,
	}
}
