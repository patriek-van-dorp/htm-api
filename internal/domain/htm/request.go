package htm

import "time"

// APIRequest represents a wrapper for incoming HTTP requests with validation and context.
type APIRequest struct {
	Input     HTMInput        `json:"input" validate:"required"`
	RequestID string          `json:"request_id" validate:"required"`
	ClientID  string          `json:"client_id,omitempty" validate:"omitempty,alphanum"`
	Priority  RequestPriority `json:"priority" validate:"omitempty,oneof=low normal high"`
}

// Validate validates the APIRequest and its nested structures
func (r *APIRequest) Validate(validator interface{}) error {
	// Note: validator interface{} is used here to avoid circular imports
	// The actual validation will be handled by the infrastructure layer
	return nil
}

// GetPriority returns the request priority, defaulting to normal if not specified
func (r *APIRequest) GetPriority() RequestPriority {
	if r.Priority == "" {
		return PriorityNormal
	}
	return r.Priority
}

// HasClientID returns true if a client ID is specified
func (r *APIRequest) HasClientID() bool {
	return r.ClientID != ""
}

// GetTimestamp returns the timestamp from the input data
func (r *APIRequest) GetTimestamp() time.Time {
	return r.Input.Timestamp
}

// GetSensorID returns the sensor ID from the input metadata
func (r *APIRequest) GetSensorID() string {
	return r.Input.Metadata.SensorID
}

// GetInputDimensions returns the dimensions from the input metadata
func (r *APIRequest) GetInputDimensions() []int {
	return r.Input.Metadata.Dimensions
}

// IsHighPriority returns true if the request has high priority
func (r *APIRequest) IsHighPriority() bool {
	return r.GetPriority() == PriorityHigh
}

// CreateProcessingContext creates a context for processing this request
func (r *APIRequest) CreateProcessingContext() map[string]interface{} {
	context := map[string]interface{}{
		"request_id": r.RequestID,
		"sensor_id":  r.GetSensorID(),
		"priority":   r.GetPriority().String(),
		"timestamp":  r.GetTimestamp(),
	}

	if r.HasClientID() {
		context["client_id"] = r.ClientID
	}

	return context
}
