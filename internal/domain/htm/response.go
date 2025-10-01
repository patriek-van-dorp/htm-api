package htm

import "time"

// APIResponse represents a wrapper for HTTP responses with error handling.
type APIResponse struct {
	Result       *ProcessingResult `json:"result,omitempty"`
	Error        *APIError         `json:"error,omitempty"`
	RequestID    string            `json:"request_id" validate:"required"`
	ResponseTime time.Time         `json:"response_time" validate:"required"`
}

// NewSuccessResponse creates a successful API response
func NewSuccessResponse(requestID string, result *ProcessingResult) *APIResponse {
	return &APIResponse{
		Result:       result,
		Error:        nil,
		RequestID:    requestID,
		ResponseTime: time.Now(),
	}
}

// NewErrorResponse creates an error API response
func NewErrorResponse(requestID string, apiError *APIError) *APIResponse {
	return &APIResponse{
		Result:       nil,
		Error:        apiError,
		RequestID:    requestID,
		ResponseTime: time.Now(),
	}
}

// IsSuccess returns true if the response represents a successful operation
func (r *APIResponse) IsSuccess() bool {
	return r.Error == nil && r.Result != nil
}

// IsError returns true if the response represents an error
func (r *APIResponse) IsError() bool {
	return r.Error != nil
}

// Validate ensures the response structure is valid
func (r *APIResponse) Validate() error {
	// Either Result or Error must be non-nil, but not both
	if r.Result != nil && r.Error != nil {
		return &ValidationError{
			Field:   "response",
			Message: "response cannot have both result and error",
		}
	}

	if r.Result == nil && r.Error == nil {
		return &ValidationError{
			Field:   "response",
			Message: "response must have either result or error",
		}
	}

	if r.RequestID == "" {
		return &ValidationError{
			Field:   "request_id",
			Message: "request_id is required",
		}
	}

	if r.ResponseTime.IsZero() {
		return &ValidationError{
			Field:   "response_time",
			Message: "response_time is required",
		}
	}

	return nil
}

// GetProcessingTime returns the processing time if available
func (r *APIResponse) GetProcessingTime() time.Duration {
	if r.Result != nil {
		return r.Result.Metadata.GetProcessingDuration()
	}
	return 0
}

// GetStatus returns the processing status if available
func (r *APIResponse) GetStatus() ProcessingStatus {
	if r.Result != nil {
		return r.Result.Status
	}
	if r.Error != nil {
		return StatusFailed
	}
	return ""
}
