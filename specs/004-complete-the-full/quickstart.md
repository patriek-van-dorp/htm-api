# Quickstart: Complete Spatial Pooler Engine Integration

**Date**: October 2, 2025  
**Feature**: Complete spatial pooler engine integration  
**Audience**: Developers completing HTM spatial pooler integration

## Overview

This quickstart guide demonstrates how to complete the spatial pooler engine integration, enabling the existing TDD test suite to run against the actual implementation. The integration involves wiring the production spatial pooler engine in `main.go`, configuring all dependencies, and validating that the complete system meets HTM theory requirements and performance specifications.

**Integration Architecture**: HTTP Request → Router → Handler → Service → **Actual Spatial Pooler Engine** → HTM-Compliant SDR Response

**Key Achievement**: Transform the current foundation into a production-ready HTM spatial pooler that passes all existing tests and meets performance requirements.

## Prerequisites

### System Requirements
- Go 1.23+ installed and configured
- Access to the HTM repository with spatial pooler foundation
- Understanding of HTM theory and spatial pooling principles
- Familiarity with existing test suite structure

### Dependencies Verification
```bash
# Verify Go version
go version

# Verify gonum availability
go list -m gonum.org/v1/gonum

# Check existing spatial pooler foundation
ls internal/cortical/spatial/
```

### Branch Setup
```bash
# Create integration branch
git checkout -b 004-complete-the-full

# Verify existing foundation
git log --oneline | grep "spatial.*pooler"
```

## Quick Start Steps

### Step 1: Verify Current Implementation Status

**Objective**: Understand what exists and what needs completion

```bash
# Check spatial pooler implementation
find . -name "*.go" -path "*/spatial/*" -exec basename {} \;

# Verify test suite structure
find tests/ -name "*spatial*" -type f

# Check current main.go integration
grep -n "spatial" cmd/api/main.go || echo "No spatial integration found"
```

**Expected Output**: 
- Spatial pooler components exist in `internal/cortical/spatial/`
- Test suite has comprehensive spatial pooler tests
- `main.go` lacks complete integration

### Step 2: Complete Main.go Integration

**Objective**: Wire all components together in the application entry point

```bash
# Backup current main.go
cp cmd/api/main.go cmd/api/main.go.backup

# Check current dependencies in main.go
grep -A 10 -B 10 "func main" cmd/api/main.go
```

**Integration Pattern**:
1. Initialize spatial pooler with configuration
2. Create spatial pooling service with actual implementation
3. Wire service into HTTP handlers
4. Configure dependency injection
5. Start server with complete integration

### Step 3: Configure Production Settings

**Objective**: Set up production-ready configuration for spatial pooler

```bash
# Check existing configuration structure
find . -name "*config*" -path "*/pkg/*" -o -path "*/internal/*"

# Verify spatial pooler configuration exists
grep -r "SpatialPoolerConfig" internal/
```

**Configuration Areas**:
- Default HTM parameters (sparsity, learning rates, etc.)
- Performance settings (timeouts, concurrency limits)
- Operational settings (logging, metrics)
- Environment-specific overrides

### Step 4: Run Integration Tests

**Objective**: Verify that existing tests pass with actual implementation

```bash
# Run spatial pooler contract tests
go test -v ./tests/contract/ -run ".*spatial.*"

# Run integration tests
go test -v ./tests/integration/ -run ".*spatial.*"

# Run performance tests
go test -v ./tests/integration/ -run ".*performance.*"
```

**Success Criteria**:
- All contract tests pass (API behavior validation)
- All integration tests pass (end-to-end validation)
- Performance tests meet <100ms requirement
- No test skips or mocks in spatial pooler path

### Step 5: Validate HTM Properties

**Objective**: Ensure implementation maintains HTM biological and mathematical properties

```bash
# Start the server
go run cmd/api/main.go

# Test HTM property validation endpoint
curl -X GET http://localhost:8080/api/v1/spatial-pooler/validation/htm-properties \
  -H "Content-Type: application/json"

# Test spatial pooling with real data
curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "encoder_output": {
        "width": 2048,
        "active_bits": [10, 25, 67, 89, 134, 256, 445, 678, 789, 901],
        "sparsity": 0.0048828125
      },
      "input_id": "test-uuid-123",
      "learning_enabled": true
    },
    "request_id": "integration-test-001",
    "client_id": "quickstart-client"
  }'
```

**Validation Checklist**:
- ✅ Sparsity level between 2-5% (HTM requirement)
- ✅ Processing time under 100ms (performance requirement)
- ✅ Semantic similarity preservation (overlap patterns)
- ✅ Learning adaptation over time
- ✅ Deterministic behavior for same inputs

### Step 6: Performance Validation

**Objective**: Confirm system meets all performance requirements

```bash
# Load test with multiple concurrent requests
for i in {1..50}; do
  curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
    -H "Content-Type: application/json" \
    -d @test_data/sample_input.json &
done
wait

# Check performance metrics
curl -X GET http://localhost:8080/api/v1/spatial-pooler/status
```

**Performance Targets**:
- Response time: <100ms per request
- Concurrent requests: Support 100 simultaneous
- Memory usage: Efficient matrix operations
- Throughput: Handle production workloads

## Common Use Cases

### Use Case 1: Development Testing
```bash
# Run specific test against actual implementation
go test -v ./tests/contract/spatial_pooler_process_test.go

# Validate test uses actual implementation (not mocks)
grep -n "mock" tests/contract/spatial_pooler_process_test.go || echo "No mocks found ✅"
```

### Use Case 2: Configuration Tuning
```bash
# Test different sparsity settings
curl -X PUT http://localhost:8080/api/v1/spatial-pooler/config \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "sparsity_ratio": 0.025,
      "learning_enabled": true,
      "performance_mode": "optimized"
    }
  }'

# Validate new configuration
curl -X GET http://localhost:8080/api/v1/spatial-pooler/config
```

### Use Case 3: Production Deployment Readiness
```bash
# Health check validation
curl -X GET http://localhost:8080/api/v1/health

# Verify all components show "production" implementation
curl -X GET http://localhost:8080/api/v1/spatial-pooler/status | \
  jq '.engine_status.implementation'
```

## Integration Examples

### Complete Main.go Integration Pattern

```go
// Example integration structure (not actual code)
func main() {
    // 1. Load configuration
    config := loadConfiguration()
    
    // 2. Initialize spatial pooler engine
    spatialPooler, err := spatial.NewSpatialPooler(config.SpatialPooler)
    if err != nil {
        log.Fatal("Failed to initialize spatial pooler:", err)
    }
    
    // 3. Create service with actual implementation
    spatialService := services.NewSpatialPoolingService(spatialPooler)
    
    // 4. Wire into HTTP handlers
    handler := handlers.NewSpatialPoolerHandler(spatialService)
    
    // 5. Start server
    router := setupRouter(handler)
    log.Fatal(router.Run(":8080"))
}
```

### Service Integration Pattern

```go
// Example service wiring (not actual code)
type SpatialPoolingServiceImpl struct {
    pooler *spatial.SpatialPooler  // Actual implementation
    logger Logger
    metrics MetricsCollector
}

func (s *SpatialPoolingServiceImpl) ProcessSpatialPooling(
    ctx context.Context, 
    input *htm.PoolingInput,
) (*htm.PoolingResult, error) {
    // Use actual spatial pooler, not mock
    return s.pooler.Process(input)
}
```

## Testing and Validation

### Comprehensive Test Execution

```bash
# Run all spatial pooler tests with actual implementation
make test-spatial-pooler-integration

# Run performance benchmarks
make benchmark-spatial-pooler

# Run HTM property validation
make validate-htm-properties
```

### Manual Validation Steps

1. **Sparsity Validation**:
   ```bash
   # Submit test input and verify output sparsity
   curl -X POST .../spatial-pooler/process -d @test_input.json | \
     jq '.result.sparsity_level'
   # Should be between 0.02 and 0.05
   ```

2. **Performance Validation**:
   ```bash
   # Check processing time
   curl -X POST .../spatial-pooler/process -d @test_input.json | \
     jq '.result.processing_time_ms'
   # Should be < 100ms
   ```

3. **Learning Validation**:
   ```bash
   # Submit same input twice, verify adaptation
   curl -X POST .../spatial-pooler/process -d @same_input.json
   curl -X POST .../spatial-pooler/process -d @same_input.json
   # Compare outputs for learning effects
   ```

## Troubleshooting

### Common Integration Issues

1. **Spatial Pooler Not Initialized**
   ```bash
   # Check error logs
   grep -i "spatial.*pooler.*init" logs/application.log
   
   # Verify configuration
   curl -X GET http://localhost:8080/api/v1/spatial-pooler/config
   ```

2. **Tests Still Using Mocks**
   ```bash
   # Search for mock usage in tests
   grep -r "mock" tests/ | grep -i spatial
   
   # Verify actual implementation is wired
   grep -r "NewSpatialPooler" cmd/api/main.go
   ```

3. **Performance Issues**
   ```bash
   # Check processing times
   curl -X GET http://localhost:8080/api/v1/spatial-pooler/status | \
     jq '.runtime_metrics.average_processing_time_ms'
   
   # Monitor memory usage
   curl -X GET http://localhost:8080/api/v1/health | \
     jq '.components.memory'
   ```

4. **HTM Property Violations**
   ```bash
   # Validate HTM properties
   curl -X GET http://localhost:8080/api/v1/spatial-pooler/validation/htm-properties
   
   # Check sparsity levels
   curl -X POST .../spatial-pooler/process -d @test.json | \
     jq '.result.sparsity_level'
   ```

### Performance Optimization

1. **Matrix Operation Optimization**:
   - Verify gonum is used for all matrix operations
   - Check memory allocation patterns
   - Monitor garbage collection impact

2. **Concurrency Optimization**:
   - Test concurrent request handling
   - Verify thread safety
   - Monitor resource contention

3. **Memory Optimization**:
   - Monitor heap size and growth
   - Check for memory leaks
   - Optimize matrix reuse

## Next Steps

### Immediate Validation
1. **Complete Integration**: Verify all components wired in main.go
2. **Test Suite Validation**: Ensure all tests pass with actual implementation
3. **Performance Validation**: Confirm <100ms processing time requirement
4. **HTM Validation**: Verify biological and mathematical properties

### Future Enhancements
1. **Temporal Memory Integration**: Prepare for next cortical column component
2. **Advanced Configuration**: Add more sophisticated parameter tuning
3. **Monitoring Enhancement**: Expand operational metrics and alerting
4. **Optimization**: Further performance improvements based on production load

### Production Readiness
1. **Load Testing**: Validate system under production-like load
2. **Error Handling**: Comprehensive error scenarios testing
3. **Monitoring Setup**: Operational dashboards and alerting
4. **Documentation**: Complete API documentation and operational guides

## Support

- **Integration Documentation**: See `/contracts/integration-api-contracts.md` for API details
- **Performance Benchmarks**: Monitor `/api/v1/spatial-pooler/status` for metrics
- **HTM Validation**: Use `/api/v1/spatial-pooler/validation/htm-properties` for property validation
- **Error Debugging**: Check `/api/v1/health` for system status and component health

This quickstart provides the foundation for completing the spatial pooler engine integration and enabling comprehensive test validation against the actual HTM implementation. The integration serves as the cornerstone for future temporal memory components and production HTM deployments.