package sensors

import (
	"errors"
	"fmt"

	"github.com/htm-project/neural-api/internal/sensors/sdr"
)

// SDRValidator provides validation logic for SDRs and encoding operations
type SDRValidator struct {
	sparsityManager    *sdr.SparsityManager
	inputSizeValidator *InputSizeValidator
	silentFailureMode  bool
}

// NewSDRValidator creates a new SDR validator with specified target sparsity
func NewSDRValidator(targetSparsity float64, silentFailureMode bool) (*SDRValidator, error) {
	sparsityManager, err := sdr.NewSparsityManager(targetSparsity)
	if err != nil {
		return nil, err
	}

	return &SDRValidator{
		sparsityManager:    sparsityManager,
		inputSizeValidator: NewInputSizeValidator(),
		silentFailureMode:  silentFailureMode,
	}, nil
}

// ValidateSDR checks if an SDR meets all HTM compliance requirements
func (v *SDRValidator) ValidateSDR(sdr *sdr.SDR) error {
	if sdr == nil {
		return errors.New("SDR cannot be nil")
	}

	// Check HTM sparsity compliance
	if err := v.sparsityManager.ValidateSDRSparsity(sdr); err != nil {
		return &ValidationError{
			Component: "sparsity",
			Reason:    err.Error(),
		}
	}

	// Check basic SDR properties
	if sdr.Width() <= 0 {
		return &ValidationError{
			Component: "width",
			Reason:    "SDR width must be positive",
		}
	}

	activeBits := sdr.ActiveBits()
	if len(activeBits) > sdr.Width() {
		return &ValidationError{
			Component: "active_bits",
			Reason:    "number of active bits cannot exceed SDR width",
		}
	}

	// Check that active bits are in valid range and sorted
	for i, bit := range activeBits {
		if bit < 0 || bit >= sdr.Width() {
			return &ValidationError{
				Component: "active_bits",
				Reason:    fmt.Sprintf("active bit %d out of range [0, %d)", bit, sdr.Width()),
			}
		}

		if i > 0 && activeBits[i-1] >= bit {
			return &ValidationError{
				Component: "active_bits",
				Reason:    "active bits must be sorted and unique",
			}
		}
	}

	return nil
}

// ValidateInput checks if input is valid for encoding
func (v *SDRValidator) ValidateInput(input interface{}) error {
	if input == nil {
		return &ValidationError{
			Component: "input",
			Reason:    "input cannot be nil",
		}
	}

	// Check input size limit
	if err := v.inputSizeValidator.ValidateInputSize(input); err != nil {
		if v.silentFailureMode {
			// In silent failure mode, oversized inputs are handled gracefully
			return nil
		}
		return err
	}

	return nil
}

// ShouldTriggerSilentFailure determines if input should trigger silent failure
func (v *SDRValidator) ShouldTriggerSilentFailure(input interface{}) bool {
	if !v.silentFailureMode {
		return false
	}

	// Check input size limit
	if v.inputSizeValidator.ShouldTriggerSilentFailure(input) {
		return true
	}

	// Add other silent failure triggers here (e.g., invalid input types)
	return false
}

// CreateEmptySDR creates an empty SDR for silent failure scenarios
func (v *SDRValidator) CreateEmptySDR(width int) (*sdr.SDR, error) {
	return sdr.NewEmptySDR(width)
}

// ValidateConfiguration checks if sensor configuration is valid
func (v *SDRValidator) ValidateConfiguration(config *SensorConfig) error {
	if config == nil {
		return &ValidationError{
			Component: "configuration",
			Reason:    "configuration cannot be nil",
		}
	}

	return config.IsValid()
}

// EncodingValidator provides validation for the encoding process
type EncodingValidator struct {
	sdrValidator *SDRValidator
}

// NewEncodingValidator creates a new encoding validator
func NewEncodingValidator(targetSparsity float64, silentFailureMode bool) (*EncodingValidator, error) {
	sdrValidator, err := NewSDRValidator(targetSparsity, silentFailureMode)
	if err != nil {
		return nil, err
	}

	return &EncodingValidator{
		sdrValidator: sdrValidator,
	}, nil
}

// ValidateEncodingOperation performs comprehensive validation of an encoding operation
func (v *EncodingValidator) ValidateEncodingOperation(input interface{}, config *SensorConfig, output *sdr.SDR) error {
	// Validate input
	if err := v.sdrValidator.ValidateInput(input); err != nil {
		return err
	}

	// Validate configuration
	if err := v.sdrValidator.ValidateConfiguration(config); err != nil {
		return err
	}

	// Validate output SDR
	if err := v.sdrValidator.ValidateSDR(output); err != nil {
		return err
	}

	// Check that output SDR width matches configuration
	if output.Width() != config.SDRWidth {
		return &ValidationError{
			Component: "output_sdr",
			Reason: fmt.Sprintf("output SDR width %d does not match configured width %d",
				output.Width(), config.SDRWidth),
		}
	}

	// Check that output sparsity is close to target
	expectedSparsity := config.TargetSparsity
	actualSparsity := output.Sparsity()
	tolerance := 0.01 // 1% tolerance

	if actualSparsity < expectedSparsity-tolerance || actualSparsity > expectedSparsity+tolerance {
		return &ValidationError{
			Component: "output_sparsity",
			Reason: fmt.Sprintf("output sparsity %.3f differs from target %.3f by more than tolerance %.3f",
				actualSparsity, expectedSparsity, tolerance),
		}
	}

	return nil
}

// SensorValidator provides validation specific to sensor implementations
type SensorValidator struct {
	encodingValidator *EncodingValidator
}

// NewSensorValidator creates a new sensor validator
func NewSensorValidator(targetSparsity float64, silentFailureMode bool) (*SensorValidator, error) {
	encodingValidator, err := NewEncodingValidator(targetSparsity, silentFailureMode)
	if err != nil {
		return nil, err
	}

	return &SensorValidator{
		encodingValidator: encodingValidator,
	}, nil
}

// ValidateSensorImplementation checks if a sensor correctly implements the interface
func (v *SensorValidator) ValidateSensorImplementation(sensor SensorInterface, testInputs []interface{}) error {
	if sensor == nil {
		return &ValidationError{
			Component: "sensor",
			Reason:    "sensor cannot be nil",
		}
	}

	// Check that sensor has valid metadata
	metadata := sensor.Metadata()
	if metadata.Type == "" {
		return &ValidationError{
			Component: "metadata",
			Reason:    "sensor type cannot be empty",
		}
	}

	if metadata.SDRWidth <= 0 {
		return &ValidationError{
			Component: "metadata",
			Reason:    "SDR width in metadata must be positive",
		}
	}

	if metadata.MaxInputSize <= 0 {
		return &ValidationError{
			Component: "metadata",
			Reason:    "max input size in metadata must be positive",
		}
	}

	// Test encoding with provided inputs
	for i, input := range testInputs {
		sdr, err := sensor.Encode(input)
		if err != nil {
			return &ValidationError{
				Component: "encoding",
				Reason:    fmt.Sprintf("encoding failed for test input %d: %v", i, err),
			}
		}

		if sdr == nil {
			return &ValidationError{
				Component: "encoding",
				Reason:    fmt.Sprintf("encoding returned nil SDR for test input %d", i),
			}
		}

		// Check consistency: same input should produce same output
		sdr2, err := sensor.Encode(input)
		if err != nil {
			return &ValidationError{
				Component: "consistency",
				Reason:    fmt.Sprintf("second encoding failed for test input %d: %v", i, err),
			}
		}

		if len(sdr.ActiveBits()) != len(sdr2.ActiveBits()) {
			return &ValidationError{
				Component: "consistency",
				Reason:    fmt.Sprintf("inconsistent encoding for test input %d", i),
			}
		}

		// Check that active bits are identical
		activeBits1 := sdr.ActiveBits()
		activeBits2 := sdr2.ActiveBits()
		for j := range activeBits1 {
			if activeBits1[j] != activeBits2[j] {
				return &ValidationError{
					Component: "consistency",
					Reason:    fmt.Sprintf("inconsistent active bits for test input %d", i),
				}
			}
		}
	}

	return nil
}

// PerformanceValidator validates that encoding operations meet performance requirements
type PerformanceValidator struct {
	maxEncodingTime float64 // Maximum encoding time in milliseconds
}

// NewPerformanceValidator creates a performance validator with sub-millisecond requirement
func NewPerformanceValidator() *PerformanceValidator {
	return &PerformanceValidator{
		maxEncodingTime: 1.0, // 1 millisecond maximum
	}
}

// GetMaxEncodingTime returns the maximum allowed encoding time
func (v *PerformanceValidator) GetMaxEncodingTime() float64 {
	return v.maxEncodingTime
}

// SetMaxEncodingTime sets a custom maximum encoding time
func (v *PerformanceValidator) SetMaxEncodingTime(maxTimeMs float64) {
	v.maxEncodingTime = maxTimeMs
}
