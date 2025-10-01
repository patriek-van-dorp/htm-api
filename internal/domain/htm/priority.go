package htm

import "fmt"

// RequestPriority represents the processing priority level for requests.
type RequestPriority string

const (
	PriorityLow    RequestPriority = "low"
	PriorityNormal RequestPriority = "normal"
	PriorityHigh   RequestPriority = "high"
)

// IsValid returns true if the priority is a valid RequestPriority value
func (p RequestPriority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityNormal, PriorityHigh:
		return true
	default:
		return false
	}
}

// String returns the string representation of the priority
func (p RequestPriority) String() string {
	return string(p)
}

// GetWeight returns a numeric weight for the priority (higher number = higher priority)
func (p RequestPriority) GetWeight() int {
	switch p {
	case PriorityHigh:
		return 3
	case PriorityNormal:
		return 2
	case PriorityLow:
		return 1
	default:
		return 2 // Default to normal
	}
}

// IsHigherThan returns true if this priority is higher than the other priority
func (p RequestPriority) IsHigherThan(other RequestPriority) bool {
	return p.GetWeight() > other.GetWeight()
}

// MarshalText implements encoding.TextMarshaler for JSON serialization
func (p RequestPriority) MarshalText() ([]byte, error) {
	if !p.IsValid() {
		return nil, fmt.Errorf("invalid RequestPriority: %s", string(p))
	}
	return []byte(p), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for JSON deserialization
func (p *RequestPriority) UnmarshalText(text []byte) error {
	priority := RequestPriority(text)
	if !priority.IsValid() {
		return fmt.Errorf("invalid RequestPriority: %s", string(text))
	}
	*p = priority
	return nil
}

// AllPriorities returns all valid RequestPriority values
func AllPriorities() []RequestPriority {
	return []RequestPriority{
		PriorityLow,
		PriorityNormal,
		PriorityHigh,
	}
}

// PriorityFromString converts a string to RequestPriority with validation
func PriorityFromString(s string) (RequestPriority, error) {
	priority := RequestPriority(s)
	if !priority.IsValid() {
		return "", fmt.Errorf("invalid RequestPriority: %s", s)
	}
	return priority, nil
}

// GetDefaultPriority returns the default priority level
func GetDefaultPriority() RequestPriority {
	return PriorityNormal
}
