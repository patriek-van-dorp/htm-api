# Data Model: Generic HTM Sensor Package

## Core Entities

### SDR (Sparsely Distributed Representation)
**Purpose**: Binary vector representation with sparse activation pattern
**Key Attributes**:
- `Width`: Total number of bits in the representation (e.g., 2048)
- `ActiveBits`: Slice of indices indicating which bits are active
- `Sparsity`: Percentage of active bits (typically 2-5%)
- `Timestamp`: Optional temporal marker for sequence learning

**Validation Rules**:
- Width must be > 0 and typically 1000-10000 bits
- ActiveBits must be unique, sorted indices within [0, Width)
- Sparsity must be between 1-10% for biological plausibility
- ActiveBits length must match calculated sparsity

**State Transitions**:
- Immutable once created for thread safety
- Builder pattern for construction with validation
- Comparison operations for overlap and similarity

### SensorInterface
**Purpose**: Contract defining how any sensor converts input to SDR
**Key Methods**:
- `Encode(input interface{}) (SDR, error)`: Primary encoding function
- `Configure(params SensorConfig) error`: Set encoding parameters
- `Validate() error`: Check sensor configuration validity
- `Metadata() SensorMetadata`: Get sensor characteristics

**Validation Rules**:
- Encode must return consistent SDR for same input
- Configure must validate parameter ranges and compatibility
- Metadata must accurately reflect encoding capabilities

### SensorConfig
**Purpose**: Configuration parameters for sensor encoding behavior
**Key Attributes**:
- `SDRWidth`: Output SDR bit width
- `TargetSparsity`: Desired active bit percentage
- `Resolution`: Encoding precision (for numeric sensors)
- `Range`: Valid input range (for bounded sensors)
- `CustomParams`: Type-specific configuration map

**Validation Rules**:
- SDRWidth must be positive and power-of-2 preferred
- TargetSparsity must be 0.01-0.10 (1-10%)
- Resolution must be positive for numeric encoders
- Range bounds must be valid for the data type

### NumericEncoder
**Purpose**: Encodes scalar numeric values into SDRs
**Key Attributes**:
- `MinValue`, `MaxValue`: Input range boundaries
- `Resolution`: Smallest distinguishable difference
- `Periodic`: Whether values wrap around (e.g., angles)
- `ClipInput`: Whether to clip out-of-range values

**Validation Rules**:
- MinValue < MaxValue for valid range
- Resolution > 0 and < (MaxValue - MinValue)
- Periodic only valid for bounded ranges
- SDR consistency for values within resolution

### CategoricalEncoder
**Purpose**: Encodes discrete categories into SDRs
**Key Attributes**:
- `Categories`: Known category set
- `HashFunction`: Collision handling for large vocabularies
- `BucketSize`: SDR bits per category
- `SimilarityMapping`: Optional semantic relationships

**Validation Rules**:
- Categories must be non-empty and unique
- BucketSize must allow for target sparsity
- HashFunction must provide consistent mapping
- SimilarityMapping preserves semantic distance

### TextEncoder
**Purpose**: Encodes text strings into SDRs with semantic meaning
**Key Attributes**:
- `Tokenizer`: Text splitting strategy
- `VocabularySize`: Maximum unique tokens
- `ContextWindow`: Token sequence length
- `SemanticSimilarity`: Preservation of word relationships

**Validation Rules**:
- Tokenizer must handle edge cases (empty, special chars)
- VocabularySize must fit in available SDR space
- ContextWindow must be positive and reasonable
- Output SDRs preserve semantic similarity

### SpatialEncoder
**Purpose**: Encodes 2D spatial data (images, coordinates) into SDRs
**Key Attributes**:
- `InputDimensions`: Width × Height of input space
- `ReceptiveFieldSize`: Local encoding window
- `Overlap`: Receptive field overlap percentage
- `TopologyPreservation`: Maintain spatial relationships

**Validation Rules**:
- InputDimensions must be positive
- ReceptiveFieldSize must be < InputDimensions
- Overlap must be 0-0.9 for effective encoding
- Topology preservation verified through neighbor similarity

### SensorRegistry
**Purpose**: Management system for discovering and instantiating sensors
**Key Attributes**:
- `RegisteredSensors`: Map of sensor type to factory function
- `DefaultConfigs`: Standard configurations per sensor type
- `ValidationRules`: Type-specific validation logic

**Validation Rules**:
- Sensor types must be unique identifiers
- Factory functions must return valid sensor instances
- Default configs must pass sensor validation
- Registration must be thread-safe

### ValidationMetrics
**Purpose**: Quality measures for SDR outputs and sensor performance
**Key Attributes**:
- `SparsityMeasure`: Actual vs target sparsity
- `OverlapAnalysis`: Similarity distribution for similar inputs
- `ConsistencyCheck`: Same input produces same SDR
- `SemanticPreservation`: Distance relationships maintained

**Validation Rules**:
- Sparsity tolerance must be within ±0.5% of target
- Overlap must increase with input similarity
- Consistency requires 100% bit match for identical inputs
- Semantic preservation within configurable tolerance

### BatchProcessor
**Purpose**: Handles batch processing of multiple inputs for efficiency
**Key Attributes**:
- `InputBatch`: Collection of input data items
- `BatchSize`: Maximum number of items per batch
- `ProcessingMode`: Sequential or parallel processing strategy
- `OutputCollection`: Resulting SDR collection

**Validation Rules**:
- BatchSize must be positive and within memory limits
- All inputs in batch must be compatible with sensor type
- ProcessingMode must respect single-threaded constraint
- Output collection maintains input order correlation

### SensorComposition
**Purpose**: Combines multiple sensors for multi-modal input processing
**Key Attributes**:
- `ComponentSensors`: Collection of individual sensor instances
- `CombinationStrategy`: How to merge multiple SDRs
- `WeightingScheme`: Relative importance of each sensor
- `OutputDimensions`: Final combined SDR dimensions

**Validation Rules**:
- ComponentSensors must be compatible (similar output dimensions)
- CombinationStrategy must preserve HTM properties
- WeightingScheme must sum to 1.0 for normalized combination
- Output maintains sparsity requirements

### MultiSDRCollection
**Purpose**: Manages collections of SDRs from spatial subdivision or multi-modal encoding
**Key Attributes**:
- `SDRArray`: Ordered collection of related SDRs
- `SpatialMapping`: Coordinate system for spatial relationships
- `TemporalSequence`: Time-ordered SDR relationships
- `MetadataMap`: Additional information per SDR

**Validation Rules**:
- SDRArray must contain valid SDR instances
- SpatialMapping must preserve topology when applicable
- TemporalSequence must maintain chronological order
- All SDRs must have compatible dimensions for aggregation

## Entity Relationships

```
SensorInterface
    ↓ implements
    ├── NumericEncoder
    ├── CategoricalEncoder
    ├── TextEncoder
    └── SpatialEncoder
    
SensorConfig
    ↓ configures
    SensorInterface
    
SensorInterface
    ↓ produces
    SDR
    
SensorRegistry
    ↓ manages
    SensorInterface
    
BatchProcessor
    ↓ processes multiple
    SensorInterface → SDR[]
    
SensorComposition
    ↓ combines
    SensorInterface[] → SDR
    
SpatialEncoder
    ↓ can produce
    MultiSDRCollection
    
ValidationMetrics
    ↓ validates
    SDR + SensorInterface + BatchProcessor + SensorComposition
```

## Data Flow

1. **Registration**: Sensor types register with SensorRegistry
2. **Configuration**: SensorConfig parameters set for specific use case
3. **Instantiation**: Registry creates configured sensor instance
4. **Encoding**: Input data converted to SDR via SensorInterface.Encode()
5. **Batch Processing**: Multiple inputs processed via BatchProcessor for efficiency
6. **Sensor Composition**: Multi-modal inputs combined via SensorComposition
7. **Multi-SDR Handling**: Spatial subdivision creates MultiSDRCollection
8. **Validation**: ValidationMetrics verify SDR quality and properties
9. **Utilization**: SDR(s) used by downstream HTM algorithms

## Persistence

**Note**: This is a stateless encoding library. No persistent storage required.
- Configurations may be serialized for reproducibility
- SDRs are transient data structures for processing
- Sensor state is configuration-only, no learned parameters