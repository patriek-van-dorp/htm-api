package htm

// InputMetadata provides context and configuration for HTM input processing.
type InputMetadata struct {
	Dimensions      []int                  `json:"dimensions" validate:"required,min=2"`
	SensorID        string                 `json:"sensor_id" validate:"required,alphanum"`
	ProcessingHints map[string]interface{} `json:"processing_hints,omitempty"`
	Version         string                 `json:"version" validate:"required,oneof=v1.0"`
}

// GetMatrixShape returns the shape of the expected matrix
func (m *InputMetadata) GetMatrixShape() (rows, cols int) {
	if len(m.Dimensions) >= 2 {
		return m.Dimensions[0], m.Dimensions[1]
	}
	return 0, 0
}

// IsValidShape validates that the dimensions are positive
func (m *InputMetadata) IsValidShape() bool {
	if len(m.Dimensions) < 2 {
		return false
	}

	for _, dim := range m.Dimensions {
		if dim <= 0 {
			return false
		}
	}

	return true
}

// GetProcessingHint retrieves a processing hint by key
func (m *InputMetadata) GetProcessingHint(key string) (interface{}, bool) {
	if m.ProcessingHints == nil {
		return nil, false
	}

	value, exists := m.ProcessingHints[key]
	return value, exists
}

// SetProcessingHint sets a processing hint
func (m *InputMetadata) SetProcessingHint(key string, value interface{}) {
	if m.ProcessingHints == nil {
		m.ProcessingHints = make(map[string]interface{})
	}

	m.ProcessingHints[key] = value
}
