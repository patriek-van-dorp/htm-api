package sensors

import (
	"github.com/htm-project/neural-api/internal/sensors/sdr"
)

// SensorInterface defines the contract for all HTM sensor implementations
type SensorInterface interface {
	// Encode converts input data into a Sparsely Distributed Representation
	// Returns error if input is invalid or encoding fails
	// In silent failure mode, returns empty SDR with no error for invalid inputs
	Encode(input interface{}) (SDR, error)

	// Configure sets encoding parameters and validates configuration
	// Must be called before first Encode operation
	Configure(config SensorConfig) error

	// Validate checks if sensor configuration is valid
	// Returns detailed error information for invalid configurations
	Validate() error

	// Metadata returns sensor characteristics and capabilities
	// Used for introspection and compatibility checking
	Metadata() SensorMetadata

	// Clone creates a new sensor instance with same configuration
	// Useful for concurrent processing with shared configuration
	Clone() SensorInterface
}

// SDR interface wraps the internal SDR implementation for public API
type SDR interface {
	// Width returns the total number of bits in the representation
	Width() int

	// ActiveBits returns indices of active (1) bits in sorted order
	ActiveBits() []int

	// Sparsity returns the percentage of active bits (0.0-1.0)
	Sparsity() float64

	// IsActive returns true if the bit at given index is active
	IsActive(index int) bool

	// Overlap calculates the number of shared active bits with another SDR
	Overlap(other SDR) int

	// Similarity returns normalized overlap (0.0-1.0) with another SDR
	Similarity(other SDR) float64

	// String returns a string representation for debugging
	String() string
}

// SensorMetadata provides information about sensor characteristics
type SensorMetadata struct {
	Type         string                 // Sensor type identifier ("numeric", "categorical", etc.)
	SDRWidth     int                    // Configured SDR width
	Sparsity     float64                // Target sparsity
	MaxInputSize int                    // Maximum input size in bytes (1MB limit)
	Capabilities map[string]interface{} // Type-specific capabilities
}

// SensorFactory is a function type for creating sensor instances
type SensorFactory func() SensorInterface

// EncodingError represents an error during encoding operation
type EncodingError struct {
	SensorType string
	Input      interface{}
	Reason     string
	SilentMode bool
}

func (e *EncodingError) Error() string {
	if e.SilentMode {
		return "encoding error (silent mode): " + e.Reason
	}
	return "encoding error: " + e.Reason
}

// ConfigurationError represents an error during sensor configuration
type ConfigurationError struct {
	Parameter string
	Value     interface{}
	Reason    string
}

func (e *ConfigurationError) Error() string {
	return "configuration error for " + e.Parameter + ": " + e.Reason
}

// ValidationError represents an error during sensor validation
type ValidationError struct {
	Component string
	Reason    string
}

func (e *ValidationError) Error() string {
	return "validation error in " + e.Component + ": " + e.Reason
}

// SDRWrapper wraps the internal SDR implementation to conform to the public interface
type SDRWrapper struct {
	internal *sdr.SDR
}

// NewSDRWrapper creates a wrapper for internal SDR
func NewSDRWrapper(internalSDR *sdr.SDR) SDR {
	return &SDRWrapper{internal: internalSDR}
}

// Width returns the total number of bits in the representation
func (w *SDRWrapper) Width() int {
	return w.internal.Width()
}

// ActiveBits returns indices of active (1) bits in sorted order
func (w *SDRWrapper) ActiveBits() []int {
	return w.internal.ActiveBits()
}

// Sparsity returns the percentage of active bits (0.0-1.0)
func (w *SDRWrapper) Sparsity() float64 {
	return w.internal.Sparsity()
}

// IsActive returns true if the bit at given index is active
func (w *SDRWrapper) IsActive(index int) bool {
	return w.internal.IsActive(index)
}

// Overlap calculates the number of shared active bits with another SDR
func (w *SDRWrapper) Overlap(other SDR) int {
	if otherWrapper, ok := other.(*SDRWrapper); ok {
		return w.internal.Overlap(otherWrapper.internal)
	}
	return 0 // Can't compare with different SDR implementations
}

// Similarity returns normalized overlap (0.0-1.0) with another SDR
func (w *SDRWrapper) Similarity(other SDR) float64 {
	if otherWrapper, ok := other.(*SDRWrapper); ok {
		return w.internal.Similarity(otherWrapper.internal)
	}
	return 0.0 // Can't compare with different SDR implementations
}

// String returns a string representation for debugging
func (w *SDRWrapper) String() string {
	return w.internal.String()
}

// GetInternalSDR returns the internal SDR for package-internal use
func (w *SDRWrapper) GetInternalSDR() *sdr.SDR {
	return w.internal
}

// InputSizeValidator validates input size against 1MB limit
type InputSizeValidator struct {
	maxSize int // Maximum input size in bytes
}

// NewInputSizeValidator creates a validator with 1MB limit
func NewInputSizeValidator() *InputSizeValidator {
	return &InputSizeValidator{
		maxSize: 1024 * 1024, // 1MB
	}
}

// ValidateInputSize checks if input size is within limits
func (v *InputSizeValidator) ValidateInputSize(input interface{}) error {
	size := estimateInputSize(input)
	if size > v.maxSize {
		return &ValidationError{
			Component: "input_size",
			Reason:    "input size exceeds 1MB limit",
		}
	}
	return nil
}

// ShouldTriggerSilentFailure checks if input should trigger silent failure
func (v *InputSizeValidator) ShouldTriggerSilentFailure(input interface{}) bool {
	return estimateInputSize(input) > v.maxSize
}

// estimateInputSize estimates the size of input data in bytes
func estimateInputSize(input interface{}) int {
	switch v := input.(type) {
	case string:
		return len([]byte(v))
	case []byte:
		return len(v)
	case []rune:
		return len(v) * 4 // Rough estimate for UTF-8 (rune is int32)
	case []float64:
		return len(v) * 8
	case []float32:
		return len(v) * 4
	case []int:
		return len(v) * 8 // Assuming 64-bit ints
	case []int64:
		return len(v) * 8
	default:
		// Conservative estimate for unknown types
		return 100
	}
}
