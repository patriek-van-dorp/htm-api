# Research: Spatial Pooler (HTM Theory) Component

**Date**: October 1, 2025  
**Feature**: Spatial Pooler implementation for HTM API  
**Status**: Complete

## Research Tasks

### 1. HTM Spatial Pooler Algorithm Fundamentals

**Decision**: Implement HTM spatial pooler based on Numenta's research and biological cortical column principles  
**Rationale**: Spatial pooler is a core HTM algorithm that converts encoder outputs into sparse distributed representations through competitive learning and maintains semantic similarity through overlapping representations  
**Alternatives considered**:
- Simple thresholding normalization (rejected: doesn't provide learning/adaptation)  
- Random sparse projection (rejected: doesn't maintain semantic relationships)  
- k-Winner-Take-All without learning (rejected: doesn't adapt to input statistics)

**Key Algorithm Components**:
- **Potential Synapses**: Each column connects to a subset of input bits with configurable permanence values
- **Active Columns**: Top-k columns with highest overlap scores become active (typically 2-5% of total columns)
- **Learning**: Permanence values adjust based on correlation with active input patterns
- **Inhibition**: Global or local competitive inhibition ensures consistent sparsity levels
- **Boosting**: Infrequently used columns get activity boost to maintain representation balance

### 2. Integration with Existing Sensor Package

**Decision**: Accept existing sensor SDR outputs while planning future migration to proper HTM architecture  
**Rationale**: Current sensors output SDR objects (deviation from HTM theory), but spatial pooler can still provide value by normalizing sparsity and ensuring semantic continuity. Future migration can implement proper encoder → bit array → spatial pooler → SDR pipeline  
**Alternatives considered**:
- Immediate sensor refactoring (rejected: breaking changes to existing API)
- Ignore HTM theory compliance (rejected: violates constitution principles)
- Separate microservice (rejected: adds latency and complexity)

**Integration Architecture (Current)**:
- **Sensors**: Output `sensors.SDR` objects (existing behavior, deviation from HTM theory)
- **Spatial Pooler**: Accepts sensor SDRs, normalizes sparsity and semantic properties
- **API Pipeline**: HTTP Request → Sensor Encoding (to SDR) → Spatial Pooling (SDR normalization) → Response
- **Future Migration**: Sensors → bit arrays → Spatial Pooler → true SDRs (proper HTM architecture)

**HTM Theory Compliance**:
- Current: Encoder → SDR → Spatial Pooler → Normalized SDR
- Proper HTM: Encoder → Bit Array → Spatial Pooler → SDR
- Migration allows gradual transition to proper HTM architecture

### 3. Performance Optimization for Go Implementation

**Decision**: Use gonum for matrix operations with optimized bit manipulation for sparse representations  
**Rationale**: Spatial pooler involves high-dimensional matrix operations that benefit from optimized linear algebra libraries, while sparse bit operations need efficient bit manipulation  
**Alternatives considered**:
- Pure Go implementation (rejected: performance concerns for matrix operations)
- CGO with C libraries (rejected: deployment complexity)
- Separate GPU acceleration (rejected: overkill for initial implementation)

**Optimization Strategy**:
- Use `gonum.org/v1/gonum/mat` for permanence matrices and overlap calculations
- Implement custom sparse bit manipulation for active column selection
- Pre-allocate memory pools for frequently used data structures
- Benchmark against <10ms requirement with realistic input sizes

### 4. Learning Algorithm Implementation

**Decision**: Implement Hebbian-style learning with configurable permanence increment/decrement  
**Rationale**: HTM spatial pooler uses biological Hebbian learning ("cells that fire together, wire together") to adapt to input statistics over time  
**Alternatives considered**:
- Fixed connections without learning (rejected: doesn't adapt to input patterns)
- Gradient-based learning (rejected: not biologically plausible)
- Reinforcement learning (rejected: doesn't fit spatial pooler role)

**Learning Parameters**:
- `synPermInc`: Permanence increment for active synapses (typically 0.05)
- `synPermDec`: Permanence decrement for inactive synapses (typically 0.008)
- `synPermTrimThreshold`: Minimum permanence to maintain connection (typically 0.01)
- `synPermConnected`: Permanence threshold for connected synapse (typically 0.1)
- Learning can be disabled for inference-only mode

### 5. Configuration and Parameter Management

**Decision**: Create hierarchical configuration with sensible defaults and validation  
**Rationale**: Spatial pooler has many parameters that affect behavior; provide defaults based on HTM research while allowing customization  
**Alternatives considered**:
- Fixed parameters (rejected: different use cases need different configurations)
- Runtime parameter tuning (rejected: adds complexity without clear benefit)
- External configuration files (rejected: prefer programmatic configuration for API)

**Configuration Structure**:
```go
type SpatialPoolerConfig struct {
    InputDimensions    []int     // Input space dimensions
    ColumnDimensions   []int     // Column space dimensions  
    PotentialRadius    float64   // Potential synapse radius (0.0-1.0)
    PotentialPct       float64   // Percentage of potential synapses (0.0-1.0)
    GlobalInhibition   bool      // Global vs local inhibition
    LocalAreaDensity   float64   // Target density for local inhibition
    NumActiveColumns   int       // Number of active columns (global inhibition)
    StimulusThreshold  int       // Minimum input threshold for activation
    SynPermInc         float64   // Permanence increment for learning
    SynPermDec         float64   // Permanence decrement for learning
    SynPermConnected   float64   // Connected synapse threshold
    DutyCyclePeriod    int       // Period for duty cycle calculation
    BoostStrength      float64   // Boosting strength factor
    Seed              int64     // Random seed for deterministic behavior
}
```

### 6. Error Handling and Edge Cases

**Decision**: Implement robust error handling with fallback SDR generation for edge cases  
**Rationale**: Feature requirements specify handling of invalid inputs (all zeros, all ones, corrupted data) with fallback patterns  
**Alternatives considered**:
- Fail fast on invalid inputs (rejected: breaks pipeline continuity)
- Silent invalid processing (rejected: masks real problems)
- Input sanitization only (rejected: doesn't handle all edge cases)

**Edge Case Handling**:
- **All zeros input**: Generate random sparse pattern with target sparsity
- **All ones input**: Apply random inhibition to achieve target sparsity  
- **Oversized input**: Reject with clear error message per requirements
- **Undersized input**: Pad with zeros or reject based on configuration
- **No active columns**: Force minimum activation with lowest threshold

### 7. Testing Strategy for Scientific Validation

**Decision**: Multi-layered testing approach with biological property validation  
**Rationale**: HTM components must maintain specific mathematical and biological properties that require specialized testing beyond standard software testing  
**Alternatives considered**:
- Standard unit testing only (rejected: doesn't validate HTM properties)
- Manual validation (rejected: not scalable or repeatable)
- Property-based testing only (rejected: needs specific HTM scenarios)

**Testing Layers**:
1. **Unit Tests**: Individual function correctness, error handling, edge cases
2. **Property Tests**: Sparsity levels, overlap preservation, learning convergence
3. **Biological Validation**: Semantic similarity preservation, adaptation behavior
4. **Performance Tests**: Latency benchmarks, memory usage, throughput under load
5. **Integration Tests**: End-to-end pipeline with real sensor data

**HTM-Specific Validations**:
- Sparsity levels remain within 2-5% target range
- Semantic similarity preserved (similar inputs → overlapping SDRs)
- Learning improves representation quality over time
- Boosting maintains balanced column usage
- Deterministic mode produces identical outputs for identical inputs

## Research Conclusions

The spatial pooler implementation will follow HTM theory principles with optimizations for the Go ecosystem. Key success factors:

1. **Biological Fidelity**: Maintain HTM algorithm correctness per Numenta research
2. **Performance**: Achieve <10ms processing time through gonum optimization
3. **Integration**: Seamlessly extend existing HTM API architecture
4. **Robustness**: Handle edge cases gracefully with fallback mechanisms
5. **Testability**: Comprehensive validation of HTM properties and performance
6. **Configurability**: Flexible parameters with sensible defaults

The implementation will prepare the foundation for future temporal memory components while providing immediate value for sensor output normalization.