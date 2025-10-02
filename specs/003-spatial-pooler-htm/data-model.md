# Data Model: Spatial Pooler (HTM Theory) Component

**Date**: October 1, 2025  
**Feature**: Spatial pooler HTM implementation  
**Source**: Extracted from feature specification and research analysis

## Core Entities

### SpatialPooler
**Purpose**: Main spatial pooler algorithm implementation that transforms encoder outputs into normalized sparse distributed representations

**Fields**:
- `config SpatialPoolerConfig` - Configuration parameters for the pooler
- `permanences *mat.Dense` - Synapse permanence matrix (columns × input bits)
- `connected *mat.Dense` - Connected synapse matrix (binary, columns × input bits)  
- `activeDutyCycles []float64` - Moving average of column activation frequency
- `overlapDutyCycles []float64` - Moving average of column overlap frequency
- `boostFactors []float64` - Boost factors for infrequently used columns
- `iterationNum int64` - Current iteration number for learning
- `rng *rand.Rand` - Random number generator for deterministic/random modes

**Validation Rules**:
- Config must be valid per SpatialPoolerConfig validation
- Permanences matrix dimensions must match config input/column dimensions
- Duty cycles arrays must match number of columns
- Boost factors must be positive values
- RNG must be initialized with config seed

**State Transitions**:
- `Uninitialized` → `Configured` (via Configure())
- `Configured` → `Ready` (via internal matrix initialization)
- `Ready` → `Processing` (during Compute() execution)
- `Processing` → `Ready` (after Compute() completion)

### SpatialPoolerConfig
**Purpose**: Configuration parameters that control spatial pooler behavior and learning

**Fields**:
- `inputDimensions []int` - Input space dimensions [width, height, ...]
- `columnDimensions []int` - Column space dimensions [width, height, ...]
- `potentialRadius float64` - Radius for potential synapse connections (0.0-1.0)
- `potentialPct float64` - Percentage of potential synapses to create (0.0-1.0)
- `globalInhibition bool` - Use global vs local competitive inhibition
- `localAreaDensity float64` - Target active density for local inhibition (0.0-1.0)
- `numActiveColumns int` - Fixed number of active columns (global inhibition)
- `stimulusThreshold int` - Minimum overlap to be considered for activation
- `synPermInc float64` - Permanence increment for active synapses
- `synPermDec float64` - Permanence decrement for inactive synapses  
- `synPermConnected float64` - Threshold for connected synapse (0.0-1.0)
- `synPermTrimThreshold float64` - Minimum permanence before trimming
- `dutyCyclePeriod int` - Period for duty cycle moving average calculation
- `boostStrength float64` - Strength of boosting for inactive columns
- `seed int64` - Random seed for deterministic behavior
- `learningEnabled bool` - Enable/disable learning mode
- `mode SpatialPoolerMode` - Deterministic or randomized processing mode

**Validation Rules**:
- All dimensions must be positive integers
- Percentage values must be between 0.0 and 1.0
- numActiveColumns must be less than total number of columns (product of columnDimensions)
- stimulusThreshold must be non-negative
- dutyCyclePeriod must be positive
- synPermConnected must be between synPermTrimThreshold and 1.0
- If globalInhibition is false, localAreaDensity must be specified
- If globalInhibition is true, numActiveColumns must be specified

### SpatialPoolerMode
**Purpose**: Enumeration for spatial pooler processing modes per feature requirements

**Values**:
- `Deterministic` - Identical inputs always produce identical outputs
- `Randomized` - Controlled randomness for learning purposes

### PoolingInput
**Purpose**: Input structure for spatial pooler processing containing encoder output and metadata

**Fields**:
- `encoderOutput []byte` - Raw bit array from sensor encoder (dense/sparse bit pattern)
- `inputWidth int` - Width of the input bit array (must match spatial pooler input dimensions)
- `inputID string` - Unique identifier for tracking purposes
- `learningEnabled bool` - Override learning mode for this specific input
- `metadata map[string]interface{}` - Additional context information

**Validation Rules**:
- encoderOutput must not be nil or empty
- inputWidth must match configured spatial pooler input dimensions
- inputID must be non-empty string
- encoderOutput length must equal inputWidth (in bits)

**NOTE**: This represents the proper HTM architecture where encoders output raw bit arrays, not SDRs

### PoolingResult
**Purpose**: Output structure containing true SDR produced by spatial pooler and processing metadata

**Fields**:
- `normalizedSDR sensors.SDR` - True sparse distributed representation with consistent sparsity (2-5%)
- `inputID string` - Matching identifier from input
- `processingTime time.Duration` - Time taken for spatial pooler processing
- `activeColumns []int` - Indices of active columns in sorted order
- `avgOverlap float64` - Average overlap score for active columns
- `sparsityLevel float64` - Actual sparsity level achieved (0.0-1.0)
- `learningOccurred bool` - Whether learning modifications were applied
- `boostingApplied bool` - Whether boosting was applied to any columns

**Validation Rules**:
- normalizedSDR sparsity must be between 0.02 and 0.05 (2-5%)
- activeColumns must be sorted and within valid column range
- sparsityLevel must match normalizedSDR.Sparsity()
- processingTime should be under 10ms performance target
- activeColumns length should match numActiveColumns (global inhibition)
- Similar inputs should have 30-70% overlap in normalizedSDR active bits
- Different inputs should have <20% overlap in normalizedSDR active bits

**NOTE**: This is where true SDRs are first created in the HTM pipeline

### PoolingError
**Purpose**: Specialized error types for spatial pooler operations

**Fields**:
- `errorType PoolingErrorType` - Category of pooling error
- `message string` - Human-readable error description
- `inputID string` - Associated input identifier (if applicable)
- `configField string` - Configuration field name (for validation errors)

**Error Types**:
- `InvalidInput` - Input validation failed (oversized, corrupted, etc.)
- `ConfigurationError` - Invalid spatial pooler configuration
- `ProcessingError` - Error during spatial pooling computation
- `PerformanceError` - Processing exceeded time/memory constraints
- `LearningError` - Error during learning rule application

### SpatialPoolerMetrics
**Purpose**: Performance and behavioral metrics for monitoring and debugging

**Fields**:
- `totalProcessed int64` - Total number of inputs processed
- `averageProcessingTime time.Duration` - Moving average processing time
- `learningIterations int64` - Number of learning iterations performed
- `columnUsageDistribution []float64` - Usage frequency per column
- `averageSparsity float64` - Moving average output sparsity
- `overlapScoreDistribution []float64` - Distribution of overlap scores
- `boostingEvents int64` - Number of times boosting was applied
- `errorCounts map[PoolingErrorType]int64` - Count of errors by type

## Entity Relationships

```
SpatialPoolerConfig ──────────┐
                              │
                              ▼
SpatialPooler ◄─── PoolingInput ──► PoolingResult
     │                               │
     │                               │
     ▼                               ▼  
SpatialPoolerMetrics            sensors.SDR
     │
     ▼
PoolingError
```

**Relationship Rules**:
- SpatialPooler has one SpatialPoolerConfig (composition)
- SpatialPooler processes multiple PoolingInput instances (aggregation)
- Each PoolingInput produces one PoolingResult (1:1 transformation)
- SpatialPooler maintains one SpatialPoolerMetrics instance (composition)
- PoolingError can be associated with PoolingInput or configuration (optional association)
- PoolingResult contains normalizedSDR following existing sensor interface contracts

## Integration with Existing Types

### HTM Domain Integration
- Extends existing `internal/domain/htm` package with spatial pooler types
- Maintains compatibility with existing `htm.APIRequest` and `htm.APIResponse` patterns
- Integrates with existing `htm.HTMInput` for sensor data flow

### Sensor Package Integration (Architectural Consideration)
- **Current State**: Existing sensors output `sensors.SDR` objects (architectural deviation from HTM theory)
- **HTM Theory**: Encoders should output raw bit arrays, spatial pooler should create SDRs
- **Migration Strategy**: 
  - Phase 1: Spatial pooler accepts existing sensor SDR outputs as input
  - Phase 2 (Future): Refactor sensors to output raw bit arrays for proper HTM compliance
- **Compatibility**: Maintains full backward compatibility with existing sensor implementations
- **Separation**: Spatial pooler still handles the HTM-specific processing (sparsity normalization, semantic preservation)

### Cortical Package (New)
- **New Package**: `internal/cortical/spatial` for HTM cortical column algorithms
- **SDR Operations**: `internal/cortical/sdr` for cortical-specific SDR manipulations
- **Algorithm Focus**: Contains HTM theory implementations (spatial pooler, future temporal memory)
- **Clear Boundary**: Distinct from sensor encoding - handles post-encoding HTM processing

### API Integration
- Follows existing handler patterns in `internal/handlers`
- Uses existing validation infrastructure in `internal/infrastructure/validation`
- Maintains compatibility with existing `internal/ports` interface patterns
- **Pipeline**: HTTP Request → Sensor Encoding → Spatial Pooling → Response

### Migration Considerations

### Backward Compatibility
- All existing sensor endpoints continue to work unchanged
- New spatial pooler endpoints are additive to existing API
- Existing SDR structures remain fully compatible
- No breaking changes to existing domain types
- **Sensor Package**: Zero changes to existing sensor implementations

### Configuration Migration
- Spatial pooler configuration is separate from sensor configuration
- Existing sensor configurations remain valid and unchanged
- New cortical package configuration is independent
- Clear separation between encoding parameters (sensors) and cortical parameters (spatial pooler)

### Performance Impact
- Spatial pooler adds <10ms to processing pipeline when used
- Memory usage increases linearly with configured column dimensions
- Existing sensor-only endpoints maintain current performance characteristics
- New spatial pooler endpoints provide additional processing value
- **Optional Processing**: Spatial pooler can be bypassed for direct sensor output if needed