package htm

import "fmt"

// APIError represents structured error information for API responses.
type APIError struct {
	Code      string                 `json:"code" validate:"required"`
	Message   string                 `json:"message" validate:"required"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Retryable bool                   `json:"retryable"`
}

// Common error codes
const (
	ErrorCodeInvalidInput       = "INVALID_INPUT"
	ErrorCodeValidationError    = "VALIDATION_ERROR"
	ErrorCodeInvalidJSON        = "INVALID_JSON"
	ErrorCodeProcessingFailed   = "PROCESSING_FAILED"
	ErrorCodeTimeout            = "TIMEOUT"
	ErrorCodeInternalError      = "INTERNAL_ERROR"
	ErrorCodeUnsupportedVersion = "UNSUPPORTED_VERSION"
	ErrorCodeRateLimitExceeded  = "RATE_LIMIT_EXCEEDED"
	ErrorCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// NewAPIError creates a new APIError with the specified parameters
func NewAPIError(code, message string, retryable bool) *APIError {
	return &APIError{
		Code:      code,
		Message:   message,
		Details:   make(map[string]interface{}),
		Retryable: retryable,
	}
}

// NewValidationError creates an APIError for validation failures
func NewValidationError(message string, details map[string]interface{}) *APIError {
	return &APIError{
		Code:      ErrorCodeValidationError,
		Message:   message,
		Details:   details,
		Retryable: false,
	}
}

// NewProcessingError creates an APIError for processing failures
func NewProcessingError(message string, retryable bool) *APIError {
	return &APIError{
		Code:      ErrorCodeProcessingFailed,
		Message:   message,
		Details:   make(map[string]interface{}),
		Retryable: retryable,
	}
}

// NewInternalError creates an APIError for internal server errors
func NewInternalError(message string) *APIError {
	return &APIError{
		Code:      ErrorCodeInternalError,
		Message:   message,
		Details:   make(map[string]interface{}),
		Retryable: true,
	}
}

// NewTimeoutError creates an APIError for timeout situations
func NewTimeoutError(message string) *APIError {
	return &APIError{
		Code:      ErrorCodeTimeout,
		Message:   message,
		Details:   make(map[string]interface{}),
		Retryable: true,
	}
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// AddDetail adds a detail to the error
func (e *APIError) AddDetail(key string, value interface{}) {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
}

// GetDetail retrieves a detail by key
func (e *APIError) GetDetail(key string) (interface{}, bool) {
	if e.Details == nil {
		return nil, false
	}
	value, exists := e.Details[key]
	return value, exists
}

// IsRetryable returns true if the error indicates the request can be retried
func (e *APIError) IsRetryable() bool {
	return e.Retryable
}

// IsClientError returns true if the error is a client-side error (4xx)
func (e *APIError) IsClientError() bool {
	switch e.Code {
	case ErrorCodeInvalidInput, ErrorCodeValidationError, ErrorCodeInvalidJSON, ErrorCodeUnsupportedVersion:
		return true
	default:
		return false
	}
}

// IsServerError returns true if the error is a server-side error (5xx)
func (e *APIError) IsServerError() bool {
	switch e.Code {
	case ErrorCodeProcessingFailed, ErrorCodeTimeout, ErrorCodeInternalError, ErrorCodeServiceUnavailable:
		return true
	default:
		return false
	}
}

// GetHTTPStatusCode returns the appropriate HTTP status code for this error
func (e *APIError) GetHTTPStatusCode() int {
	switch e.Code {
	case ErrorCodeInvalidInput, ErrorCodeValidationError, ErrorCodeInvalidJSON, ErrorCodeUnsupportedVersion:
		return 400 // Bad Request
	case ErrorCodeRateLimitExceeded:
		return 429 // Too Many Requests
	case ErrorCodeProcessingFailed, ErrorCodeInternalError:
		return 500 // Internal Server Error
	case ErrorCodeTimeout:
		return 504 // Gateway Timeout
	case ErrorCodeServiceUnavailable:
		return 503 // Service Unavailable
	default:
		return 500 // Default to Internal Server Error
	}
}

// Clone creates a deep copy of the APIError
func (e *APIError) Clone() *APIError {
	details := make(map[string]interface{})
	for k, v := range e.Details {
		details[k] = v
	}

	return &APIError{
		Code:      e.Code,
		Message:   e.Message,
		Details:   details,
		Retryable: e.Retryable,
	}
}
