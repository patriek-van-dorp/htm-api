package htm

import "fmt"

// ProcessingStatus represents the status of HTM processing operations.
type ProcessingStatus string

const (
	StatusSuccess        ProcessingStatus = "SUCCESS"
	StatusPartialSuccess ProcessingStatus = "PARTIAL_SUCCESS"
	StatusFailed         ProcessingStatus = "FAILED"
	StatusTimeout        ProcessingStatus = "TIMEOUT"
	StatusRetrying       ProcessingStatus = "RETRYING"
)

// IsValid returns true if the status is a valid ProcessingStatus value
func (s ProcessingStatus) IsValid() bool {
	switch s {
	case StatusSuccess, StatusPartialSuccess, StatusFailed, StatusTimeout, StatusRetrying:
		return true
	default:
		return false
	}
}

// IsTerminal returns true if the status represents a final state
func (s ProcessingStatus) IsTerminal() bool {
	switch s {
	case StatusSuccess, StatusPartialSuccess, StatusFailed, StatusTimeout:
		return true
	case StatusRetrying:
		return false
	default:
		return false
	}
}

// IsSuccessful returns true if the status indicates successful processing
func (s ProcessingStatus) IsSuccessful() bool {
	return s == StatusSuccess || s == StatusPartialSuccess
}

// IsFailure returns true if the status indicates failed processing
func (s ProcessingStatus) IsFailure() bool {
	return s == StatusFailed || s == StatusTimeout
}

// String returns the string representation of the status
func (s ProcessingStatus) String() string {
	return string(s)
}

// MarshalText implements encoding.TextMarshaler for JSON serialization
func (s ProcessingStatus) MarshalText() ([]byte, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid ProcessingStatus: %s", string(s))
	}
	return []byte(s), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for JSON deserialization
func (s *ProcessingStatus) UnmarshalText(text []byte) error {
	status := ProcessingStatus(text)
	if !status.IsValid() {
		return fmt.Errorf("invalid ProcessingStatus: %s", string(text))
	}
	*s = status
	return nil
}

// AllStatuses returns all valid ProcessingStatus values
func AllStatuses() []ProcessingStatus {
	return []ProcessingStatus{
		StatusSuccess,
		StatusPartialSuccess,
		StatusFailed,
		StatusTimeout,
		StatusRetrying,
	}
}

// StatusFromString converts a string to ProcessingStatus with validation
func StatusFromString(s string) (ProcessingStatus, error) {
	status := ProcessingStatus(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid ProcessingStatus: %s", s)
	}
	return status, nil
}
