package sensors

import (
	"errors"
	"fmt"
)

// SensorConfig holds configuration parameters for sensor encoding behavior
type SensorConfig struct {
	SDRWidth       int                    // Output SDR bit width
	TargetSparsity float64                // Desired active bit percentage (0.01-0.10)
	Resolution     float64                // Encoding precision (for numeric sensors)
	Range          *Range                 // Valid input range (for bounded sensors)
	CustomParams   map[string]interface{} // Type-specific configuration parameters
}

// Range defines min/max bounds for numeric inputs
type Range struct {
	Min float64 // Minimum value
	Max float64 // Maximum value
}

// NewSensorConfig creates a new sensor configuration with HTM-compliant defaults
func NewSensorConfig() *SensorConfig {
	return &SensorConfig{
		SDRWidth:       2048,                         // Common HTM SDR width
		TargetSparsity: 0.02,                         // 2% sparsity (HTM recommended)
		Resolution:     0.1,                          // Default numeric resolution
		Range:          &Range{Min: 0.0, Max: 100.0}, // Default range
		CustomParams:   make(map[string]interface{}),
	}
}

// SetParam sets a custom parameter value
func (c *SensorConfig) SetParam(key string, value interface{}) {
	c.CustomParams[key] = value
}

// GetParam retrieves a custom parameter value
func (c *SensorConfig) GetParam(key string) (interface{}, bool) {
	value, exists := c.CustomParams[key]
	return value, exists
}

// GetStringParam retrieves a string parameter with default
func (c *SensorConfig) GetStringParam(key, defaultValue string) string {
	if value, exists := c.CustomParams[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

// GetIntParam retrieves an int parameter with default
func (c *SensorConfig) GetIntParam(key string, defaultValue int) int {
	if value, exists := c.CustomParams[key]; exists {
		if i, ok := value.(int); ok {
			return i
		}
	}
	return defaultValue
}

// GetFloatParam retrieves a float64 parameter with default
func (c *SensorConfig) GetFloatParam(key string, defaultValue float64) float64 {
	if value, exists := c.CustomParams[key]; exists {
		if f, ok := value.(float64); ok {
			return f
		}
	}
	return defaultValue
}

// GetBoolParam retrieves a bool parameter with default
func (c *SensorConfig) GetBoolParam(key string, defaultValue bool) bool {
	if value, exists := c.CustomParams[key]; exists {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return defaultValue
}

// ValidateSDRWidth checks if SDR width is valid
func (c *SensorConfig) ValidateSDRWidth() error {
	if c.SDRWidth <= 0 {
		return &ConfigurationError{
			Parameter: "sdr_width",
			Value:     c.SDRWidth,
			Reason:    "must be positive",
		}
	}

	// Recommend power of 2 for optimal performance, but don't enforce
	if c.SDRWidth < 100 {
		return &ConfigurationError{
			Parameter: "sdr_width",
			Value:     c.SDRWidth,
			Reason:    "recommended minimum is 100 for meaningful sparse representations",
		}
	}

	if c.SDRWidth > 100000 {
		return &ConfigurationError{
			Parameter: "sdr_width",
			Value:     c.SDRWidth,
			Reason:    "maximum recommended width is 100,000 for performance",
		}
	}

	return nil
}

// ValidateSparsity checks if target sparsity meets HTM requirements
func (c *SensorConfig) ValidateSparsity() error {
	if c.TargetSparsity < 0.01 {
		return &ConfigurationError{
			Parameter: "target_sparsity",
			Value:     c.TargetSparsity,
			Reason:    "below HTM minimum of 1% (0.01)",
		}
	}

	if c.TargetSparsity > 0.10 {
		return &ConfigurationError{
			Parameter: "target_sparsity",
			Value:     c.TargetSparsity,
			Reason:    "above HTM maximum of 10% (0.10)",
		}
	}

	return nil
}

// ValidateResolution checks if resolution is valid for numeric encoders
func (c *SensorConfig) ValidateResolution() error {
	if c.Resolution <= 0.0 {
		return &ConfigurationError{
			Parameter: "resolution",
			Value:     c.Resolution,
			Reason:    "must be positive",
		}
	}

	return nil
}

// ValidateRange checks if the input range is valid
func (c *SensorConfig) ValidateRange() error {
	if c.Range == nil {
		return nil // Range is optional
	}

	if c.Range.Min >= c.Range.Max {
		return &ConfigurationError{
			Parameter: "range",
			Value:     fmt.Sprintf("[%.3f, %.3f]", c.Range.Min, c.Range.Max),
			Reason:    "min must be less than max",
		}
	}

	return nil
}

// IsValid performs complete validation of the configuration
func (c *SensorConfig) IsValid() error {
	if err := c.ValidateSDRWidth(); err != nil {
		return err
	}

	if err := c.ValidateSparsity(); err != nil {
		return err
	}

	if err := c.ValidateResolution(); err != nil {
		return err
	}

	if err := c.ValidateRange(); err != nil {
		return err
	}

	return nil
}

// Clone creates a deep copy of the configuration
func (c *SensorConfig) Clone() *SensorConfig {
	clone := &SensorConfig{
		SDRWidth:       c.SDRWidth,
		TargetSparsity: c.TargetSparsity,
		Resolution:     c.Resolution,
		CustomParams:   make(map[string]interface{}),
	}

	// Clone range if it exists
	if c.Range != nil {
		clone.Range = &Range{
			Min: c.Range.Min,
			Max: c.Range.Max,
		}
	}

	// Clone custom parameters
	for key, value := range c.CustomParams {
		clone.CustomParams[key] = value
	}

	return clone
}

// CalculateActiveBitsCount returns expected number of active bits for this configuration
func (c *SensorConfig) CalculateActiveBitsCount() int {
	return int(float64(c.SDRWidth) * c.TargetSparsity)
}

// String returns a string representation of the configuration
func (c *SensorConfig) String() string {
	rangeStr := "nil"
	if c.Range != nil {
		rangeStr = fmt.Sprintf("[%.3f, %.3f]", c.Range.Min, c.Range.Max)
	}

	return fmt.Sprintf("SensorConfig(width=%d, sparsity=%.3f, resolution=%.3f, range=%s, params=%d)",
		c.SDRWidth, c.TargetSparsity, c.Resolution, rangeStr, len(c.CustomParams))
}

// ConfigurationTemplate provides pre-configured templates for common use cases
type ConfigurationTemplate struct {
	name   string
	config *SensorConfig
}

// PredefinedConfigurations provides common configuration templates
var PredefinedConfigurations = map[string]*ConfigurationTemplate{
	"small": {
		name: "small",
		config: &SensorConfig{
			SDRWidth:       1024,
			TargetSparsity: 0.02,
			Resolution:     0.1,
			Range:          &Range{Min: 0.0, Max: 100.0},
			CustomParams:   make(map[string]interface{}),
		},
	},
	"medium": {
		name: "medium",
		config: &SensorConfig{
			SDRWidth:       2048,
			TargetSparsity: 0.02,
			Resolution:     0.1,
			Range:          &Range{Min: 0.0, Max: 100.0},
			CustomParams:   make(map[string]interface{}),
		},
	},
	"large": {
		name: "large",
		config: &SensorConfig{
			SDRWidth:       4096,
			TargetSparsity: 0.02,
			Resolution:     0.1,
			Range:          &Range{Min: 0.0, Max: 100.0},
			CustomParams:   make(map[string]interface{}),
		},
	},
	"sparse": {
		name: "sparse",
		config: &SensorConfig{
			SDRWidth:       2048,
			TargetSparsity: 0.01,
			Resolution:     0.1,
			Range:          &Range{Min: 0.0, Max: 100.0},
			CustomParams:   make(map[string]interface{}),
		},
	},
	"dense": {
		name: "dense",
		config: &SensorConfig{
			SDRWidth:       2048,
			TargetSparsity: 0.05,
			Resolution:     0.1,
			Range:          &Range{Min: 0.0, Max: 100.0},
			CustomParams:   make(map[string]interface{}),
		},
	},
}

// GetTemplate returns a clone of the specified configuration template
func GetTemplate(name string) (*SensorConfig, error) {
	template, exists := PredefinedConfigurations[name]
	if !exists {
		return nil, errors.New("configuration template not found: " + name)
	}

	return template.config.Clone(), nil
}

// ListTemplates returns the names of available configuration templates
func ListTemplates() []string {
	names := make([]string, 0, len(PredefinedConfigurations))
	for name := range PredefinedConfigurations {
		names = append(names, name)
	}
	return names
}
