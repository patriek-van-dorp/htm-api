# Go Interface Contracts: HTM Sensor Package

## Core Interfaces

### SensorInterface
```go
// SensorInterface defines the contract for all HTM sensor implementations
type SensorInterface interface {
    // Encode converts input data into a Sparsely Distributed Representation
    // Returns error if input is invalid or encoding fails
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
```

### SDR (Sparsely Distributed Representation)
```go
// SDR represents a sparse binary vector with HTM properties
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
    
    // Validate ensures SDR meets HTM theoretical requirements
    Validate() error
    
    // String returns human-readable representation for debugging
    String() string
    
    // Bytes returns compact binary representation for serialization
    Bytes() []byte
}
```

### SensorConfig
```go
// SensorConfig contains parameters for sensor encoding behavior
type SensorConfig interface {
    // SDRWidth returns the output SDR bit width
    SDRWidth() int
    
    // TargetSparsity returns desired active bit percentage (0.0-1.0)
    TargetSparsity() float64
    
    // GetParam returns sensor-specific configuration parameter
    GetParam(key string) (interface{}, bool)
    
    // SetParam sets sensor-specific configuration parameter
    SetParam(key string, value interface{}) error
    
    // Validate ensures configuration parameters are valid
    Validate() error
    
    // Clone creates a copy of the configuration
    Clone() SensorConfig
}
```

### SensorRegistry
```go
// SensorRegistry manages sensor type registration and instantiation
type SensorRegistry interface {
    // Register adds a new sensor type with factory function
    Register(sensorType string, factory SensorFactory) error
    
    // Create instantiates a sensor of the specified type
    Create(sensorType string, config SensorConfig) (SensorInterface, error)
    
    // List returns all registered sensor types
    List() []string
    
    // IsRegistered checks if a sensor type is available
    IsRegistered(sensorType string) bool
    
    // GetMetadata returns metadata for a registered sensor type
    GetMetadata(sensorType string) (SensorTypeMetadata, error)
}
```

## Type Definitions

### SensorMetadata
```go
// SensorMetadata describes sensor capabilities and characteristics
type SensorMetadata struct {
    Type            string            // Sensor type identifier
    InputTypes      []string          // Supported input data types
    OutputSDRWidth  int               // Default SDR width
    TargetSparsity  float64           // Default sparsity level
    Description     string            // Human-readable description
    BiologicalBasis string            // Neuroscience foundation reference
    Parameters      map[string]string // Configuration parameter descriptions
}
```

### SensorTypeMetadata
```go
// SensorTypeMetadata describes a registered sensor type
type SensorTypeMetadata struct {
    Name            string              // Sensor type name
    Description     string              // Purpose and capabilities
    DefaultConfig   map[string]interface{} // Default parameter values
    RequiredParams  []string            // Mandatory configuration parameters
    OptionalParams  []string            // Optional configuration parameters
    InputTypes      []string            // Supported Go types for input
    Examples        []string            // Usage example descriptions
}
```

### SensorFactory
```go
// SensorFactory creates new sensor instances
type SensorFactory func(config SensorConfig) (SensorInterface, error)
```

### ValidationResult
```go
// ValidationResult contains SDR quality metrics
type ValidationResult struct {
    IsValid           bool              // Overall validation status
    SparsityError     float64           // Difference from target sparsity
    ConsistencyCheck  bool              // Same input produces same SDR
    SemanticCheck     bool              // Similar inputs have similar SDRs
    PerformanceMetrics map[string]float64 // Timing and memory usage
    Errors            []string          // Detailed error messages
}
```

## Extended Interfaces

### BatchProcessor
```go
// BatchProcessor handles efficient processing of multiple inputs
type BatchProcessor interface {
    // ProcessBatch encodes multiple inputs using the same sensor configuration
    // Returns slice of SDRs in same order as input slice
    ProcessBatch(sensor SensorInterface, inputs []interface{}) ([]SDR, error)
    
    // ProcessBatchWithSensors uses different sensors for different inputs
    // Length of sensors and inputs must match
    ProcessBatchWithSensors(sensors []SensorInterface, inputs []interface{}) ([]SDR, error)
    
    // Configure sets batch processing parameters
    Configure(config BatchConfig) error
    
    // GetBatchSize returns maximum items processed per batch
    GetBatchSize() int
}
```

### SensorComposition
```go
// SensorComposition combines multiple sensors for multi-modal processing
type SensorComposition interface {
    // AddSensor adds a sensor to the composition with a weight
    AddSensor(sensor SensorInterface, weight float64) error
    
    // Encode processes input through all component sensors and combines results
    Encode(input interface{}) (SDR, error)
    
    // EncodeMultiModal processes different inputs through different sensors
    EncodeMultiModal(inputs map[string]interface{}) (SDR, error)
    
    // Configure sets composition parameters
    Configure(config CompositionConfig) error
    
    // GetComponents returns the configured sensor components
    GetComponents() []SensorComponent
}
```

### MultiSDRCollection
```go
// MultiSDRCollection manages collections of related SDRs
type MultiSDRCollection interface {
    // Add appends an SDR to the collection with optional metadata
    Add(sdr SDR, metadata map[string]interface{}) error
    
    // Get retrieves SDR at specified index
    Get(index int) (SDR, error)
    
    // GetWithMetadata retrieves SDR and its metadata
    GetWithMetadata(index int) (SDR, map[string]interface{}, error)
    
    // Size returns number of SDRs in collection
    Size() int
    
    // Combine aggregates all SDRs into a single representation
    Combine(strategy CombinationStrategy) (SDR, error)
    
    // Iterator returns iterator for sequential access
    Iterator() SDRIterator
}
```

## Supporting Types

### BatchConfig
```go
// BatchConfig configures batch processing behavior
type BatchConfig struct {
    MaxBatchSize      int               // Maximum items per batch
    ProcessingMode    ProcessingMode    // Sequential or optimized processing
    MemoryLimit       int64             // Maximum memory usage in bytes
    ErrorStrategy     ErrorStrategy     // How to handle individual failures
}
```

### CompositionConfig
```go
// CompositionConfig configures sensor composition behavior
type CompositionConfig struct {
    CombinationStrategy CombinationStrategy // How to merge SDRs
    NormalizationMode   NormalizationMode   // Weight normalization approach
    OutputDimensions    int                 // Final SDR dimensions
    QualityThreshold    float64             // Minimum acceptable quality
}
```

### SensorComponent
```go
// SensorComponent represents a sensor within a composition
type SensorComponent struct {
    Sensor      SensorInterface // The sensor instance
    Weight      float64         // Relative importance (0.0-1.0)
    InputKey    string          // Key for multi-modal input mapping
    Metadata    map[string]interface{} // Additional component information
}

### SensorError
```go
// SensorError represents sensor-specific errors with context
type SensorError struct {
    SensorType  string    // Type of sensor that generated error
    Operation   string    // Operation that failed (encode, configure, validate)
    Input       string    // String representation of problematic input
    Message     string    // Human-readable error description
    Cause       error     // Underlying error if any
}

func (e *SensorError) Error() string
func (e *SensorError) Unwrap() error
```

### ConfigurationError
```go
// ConfigurationError represents invalid sensor configuration
type ConfigurationError struct {
    Parameter   string      // Parameter name that caused error
    Value       interface{} // Invalid value
    Expected    string      // Description of expected value
    Constraint  string      // Constraint that was violated
}

func (e *ConfigurationError) Error() string
```

### SDRError
```go
// SDRError represents invalid SDR properties
type SDRError struct {
    Property    string  // SDR property that failed validation
    Expected    string  // Expected value or range
    Actual      string  // Actual value that failed
    SDRWidth    int     // Width of the problematic SDR
    Sparsity    float64 // Sparsity of the problematic SDR
}

func (e *SDRError) Error() string
```

## Contract Validation Requirements

1. **Interface Compliance**: All sensor implementations must satisfy SensorInterface
2. **SDR Properties**: All SDR outputs must validate against HTM requirements
3. **Configuration Consistency**: Same configuration must produce deterministic behavior
4. **Error Handling**: All errors must use defined error types with context
5. **Thread Safety**: All interfaces must be safe for concurrent use
6. **Performance**: Encoding operations must complete within performance targets
7. **Biological Plausibility**: SDR properties must match neuroscience constraints

## Usage Contract Examples

### Basic Sensor Usage
```go
// Register sensor type
registry.Register("numeric", NewNumericSensor)

// Configure sensor
config := NewSensorConfig()
config.SetParam("sdr_width", 2048)
config.SetParam("target_sparsity", 0.02)
config.SetParam("min_value", 0.0)
config.SetParam("max_value", 100.0)

// Create and use sensor
sensor, err := registry.Create("numeric", config)
sdr, err := sensor.Encode(42.5)
```

### Multi-Sensor Pipeline
```go
// Create different sensor types
numericSensor := registry.Create("numeric", numericConfig)
textSensor := registry.Create("text", textConfig)

// Process different input types
numericSDR := numericSensor.Encode(123.45)
textSDR := textSensor.Encode("hello world")

// Validate outputs
validator := NewSDRValidator()
result := validator.Validate(numericSDR)
```

### Batch Processing
```go
// Configure batch processor
batchConfig := BatchConfig{
    MaxBatchSize: 100,
    ProcessingMode: Sequential,
    MemoryLimit: 1024 * 1024, // 1MB
}

processor := NewBatchProcessor()
processor.Configure(batchConfig)

// Process multiple inputs
inputs := []interface{}{1.2, 3.4, 5.6, 7.8}
sdrs, err := processor.ProcessBatch(numericSensor, inputs)
```

### Sensor Composition
```go
// Create composition
composition := NewSensorComposition()
composition.AddSensor(numericSensor, 0.6)
composition.AddSensor(textSensor, 0.4)

// Configure composition strategy
compConfig := CompositionConfig{
    CombinationStrategy: WeightedSum,
    OutputDimensions: 2048,
}
composition.Configure(compConfig)

// Encode multi-modal input
multiModalInput := map[string]interface{}{
    "numeric": 42.5,
    "text": "sensor data",
}
combinedSDR, err := composition.EncodeMultiModal(multiModalInput)
```

### Multi-SDR Collection
```go
// Create collection for spatial subdivision
collection := NewMultiSDRCollection()

// Add SDRs from image patches
for i, patch := range imagePatches {
    sdr, _ := spatialSensor.Encode(patch)
    metadata := map[string]interface{}{
        "patch_id": i,
        "coordinates": patch.Bounds(),
    }
    collection.Add(sdr, metadata)
}

// Combine all patches into single representation
combinedSDR, err := collection.Combine(SpatialConcatenation)
```