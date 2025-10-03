# Quickstart: Complete HTM Pipeline with Sensor-to-Motor Integration

**Feature**: Complete HTM Pipeline Implementation  
**Estimated Time**: 15-20 minutes for basic setup, 30-45 minutes for comprehensive testing  
**Prerequisites**: Go 1.23+, existing HTM API with spatial pooler

## Quick Setup Guide

### 1. Environment Verification

Verify your environment is ready for HTM pipeline development:

```bash
# Check Go version
go version  # Should be 1.23 or higher

# Verify HTM API repository
cd /path/to/htm-api
ls internal/handlers/spatial_pooler_handler.go  # Should exist

# Check existing spatial pooler functionality
go test ./tests/contract -v -run TestSpatialPooler  # Should pass

# Verify dependencies
go list -m gonum.org/v1/gonum  # Should be available
```

### 2. Branch Setup

Create and switch to the implementation branch:

```bash
# Create new feature branch
git checkout -b 006-implementation-readme-is

# Verify branch status
git status
# Should show clean working directory on new branch
```

### 3. Basic Implementation Validation

Run existing tests to establish baseline:

```bash
# Run spatial pooler tests (should pass)
go test ./tests/contract -v -run TestSpatialPoolerProcess

# Run integration tests (should pass)
go test ./tests/integration -v -run TestSpatialPooler

# Check for any existing temporal memory stubs
find . -name "*.go" -exec grep -l "temporal.*memory" {} \;
```

## Core Development Workflow

### 1. Temporal Memory Implementation

Follow test-driven development for temporal memory:

```bash
# Create temporal memory test file
touch tests/contract/temporal_memory_process_test.go

# Implement basic temporal memory test
cat > tests/contract/temporal_memory_process_test.go << 'EOF'
package contract

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestTemporalMemoryProcessBasic(t *testing.T) {
    t.Skip("Implementation pending - should fail until temporal memory is implemented")
    
    // This test should fail initially
    // Implementation will make it pass
    assert.True(t, false, "Temporal memory not yet implemented")
}
EOF

# Run test (should skip/fail)
go test ./tests/contract -v -run TestTemporalMemoryProcess
```

### 2. Motor Output Implementation

Create motor output test structure:

```bash
# Create motor output test file
touch tests/contract/motor_output_process_test.go

# Implement basic motor output test
cat > tests/contract/motor_output_process_test.go << 'EOF'
package contract

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestMotorOutputProcessBasic(t *testing.T) {
    t.Skip("Implementation pending - should fail until motor output is implemented")
    
    // This test should fail initially
    // Implementation will make it pass
    assert.True(t, false, "Motor output not yet implemented")
}
EOF

# Run test (should skip/fail)
go test ./tests/contract -v -run TestMotorOutputProcess
```

### 3. Complete Pipeline Integration

Create pipeline integration test:

```bash
# Create pipeline test file
touch tests/integration/complete_pipeline_test.go

# Implement basic pipeline test
cat > tests/integration/complete_pipeline_test.go << 'EOF'
package integration

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCompletePipelineBasic(t *testing.T) {
    t.Skip("Implementation pending - complete pipeline not yet implemented")
    
    // This test validates sensor-to-motor flow
    // Should fail until complete implementation
    assert.True(t, false, "Complete pipeline not yet implemented")
}
EOF

# Run test (should skip/fail)
go test ./tests/integration -v -run TestCompletePipeline
```

### 4. Sample Client Application

Create sample client structure:

```bash
# Create sample client directory
mkdir -p pkg/client/sample_client

# Create main client application
touch pkg/client/sample_client/main.go

# Create basic client structure
cat > pkg/client/sample_client/main.go << 'EOF'
package main

import (
    "fmt"
    "log"
)

func main() {
    fmt.Println("HTM Sample Client - Complete Pipeline Testing")
    
    // TODO: Implement sensor hosting
    // TODO: Implement HTM API integration
    // TODO: Implement motor output simulation
    
    log.Println("Sample client implementation pending")
}
EOF

# Create sensor implementations
mkdir -p pkg/client/sample_client/sensors
touch pkg/client/sample_client/sensors/temperature.go
touch pkg/client/sample_client/sensors/text.go
touch pkg/client/sample_client/sensors/image.go
touch pkg/client/sample_client/sensors/audio.go
```

## Validation Workflow

### 1. Test-Driven Development Validation

Ensure all new tests are created and initially failing:

```bash
# Run all new contract tests (should skip or fail)
go test ./tests/contract -v -run "TestTemporalMemory|TestMotorOutput"

# Run all new integration tests (should skip or fail)
go test ./tests/integration -v -run "TestCompletePipeline"

# Check test coverage baseline
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 2. API Endpoint Validation

Verify spatial pooler endpoints work before extending:

```bash
# Test existing spatial pooler endpoints
curl -X GET http://localhost:8080/api/v1/spatial-pooler/config
# Should return spatial pooler configuration

curl -X GET http://localhost:8080/api/v1/spatial-pooler/status  
# Should return spatial pooler status

# Start API server for testing
go run cmd/api/main.go &
API_PID=$!

# Wait for startup
sleep 2

# Test spatial pooler process endpoint
curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "encoder_output": {
        "width": 2048,
        "active_bits": [10, 25, 67, 89, 134]
      }
    }
  }'
# Should return spatial pooler result

# Clean up
kill $API_PID
```

### 3. HTM Compliance Validation

Check HTM biological constraints in existing implementation:

```bash
# Run HTM compliance tests
go test ./tests/contract -v -run TestHTMProperties

# Check sparsity validation
go test ./tests/contract -v -run TestSpatialPoolerSparsity

# Validate biological constraints
go test ./tests/integration -v -run TestSpatialPoolerBiological
```

## Integration Testing

### 1. Complete Pipeline Test Scenario

Create a comprehensive test scenario:

```bash
# Create test scenario file
cat > tests/integration/pipeline_scenario_test.go << 'EOF'
package integration

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCompletePipelineScenario(t *testing.T) {
    t.Skip("Implementation pending - complete scenario test")
    
    // Test scenario: Temperature sensor → HTM pipeline → Motor output
    scenario := &TestScenario{
        Name: "Basic Temperature Processing",
        Description: "Temperature sensor data flows through complete HTM pipeline",
        Duration: 30 * time.Second,
        SensorTypes: []string{"temperature"},
        ExpectedOutcomes: []string{
            "Encoder produces valid SDR",
            "Spatial pooler maintains 2-5% sparsity", 
            "Temporal memory learns sequences",
            "Motor output generates commands",
            "End-to-end latency < 100ms",
        },
    }
    
    // This test will be implemented with actual pipeline
    require.NotNil(t, scenario)
    assert.True(t, false, "Complete pipeline scenario not yet implemented")
}
EOF
```

### 2. Performance Baseline Testing

Establish performance benchmarks:

```bash
# Create performance test file
cat > tests/benchmark/pipeline_performance_test.go << 'EOF'
package benchmark

import (
    "testing"
)

func BenchmarkCompletePipeline(b *testing.B) {
    b.Skip("Implementation pending - performance benchmarks")
    
    // Benchmark complete sensor-to-motor pipeline
    // Target: <100ms end-to-end processing
    
    for i := 0; i < b.N; i++ {
        // TODO: Process sensor data through complete pipeline
        // TODO: Measure end-to-end latency
        // TODO: Validate HTM compliance
    }
}

func BenchmarkTemporalMemory(b *testing.B) {
    b.Skip("Implementation pending - temporal memory benchmarks")
    
    // Benchmark temporal memory processing alone
    // Target: <30ms processing time
    
    for i := 0; i < b.N; i++ {
        // TODO: Process spatial SDR through temporal memory
        // TODO: Measure processing latency
    }
}

func BenchmarkMotorOutput(b *testing.B) {
    b.Skip("Implementation pending - motor output benchmarks")
    
    // Benchmark motor output command generation
    // Target: <20ms command generation
    
    for i := 0; i < b.N; i++ {
        // TODO: Generate motor commands from predictions
        // TODO: Measure command generation time
    }
}
EOF

mkdir -p tests/benchmark
```

## Sample Client Testing

### 1. Basic Client Functionality

Test sample client basic operations:

```bash
# Build sample client
cd pkg/client/sample_client
go build -o sample_client main.go

# Test basic client startup
./sample_client --help
# Should show usage information (when implemented)

# Test client with mock sensors
./sample_client --sensors=temperature --duration=10s --mock-mode
# Should run mock scenario (when implemented)
```

### 2. HTM API Integration Testing

Test client integration with HTM API:

```bash
# Start HTM API server
cd ../../../
go run cmd/api/main.go &
API_PID=$!

# Test client API connectivity
cd pkg/client/sample_client
./sample_client --test-connection --api-url=http://localhost:8080
# Should verify API connectivity (when implemented)

# Run basic integration test
./sample_client --sensors=temperature --duration=30s --api-url=http://localhost:8080
# Should process temperature data through HTM pipeline (when implemented)

# Clean up
kill $API_PID
```

## Success Criteria Validation

### 1. Implementation Completeness

Verify all components are implemented:

```bash
# Check temporal memory implementation
ls internal/cortical/temporal/
# Should contain temporal memory implementation files

# Check motor output implementation  
ls internal/handlers/motor_output_handler.go
# Should exist

# Check pipeline orchestration
ls internal/handlers/pipeline_handler.go
# Should exist

# Check sample client completeness
ls pkg/client/sample_client/sensors/
# Should contain all four sensor implementations
```

### 2. Test Coverage Validation

Ensure comprehensive test coverage:

```bash
# Run all tests
go test ./... -v

# Check test coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
# Should show >80% coverage

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # Review coverage areas
```

### 3. Performance Requirements

Validate performance targets:

```bash
# Run performance benchmarks
go test -bench=. ./tests/benchmark/

# Check end-to-end latency
go test -bench=BenchmarkCompletePipeline ./tests/benchmark/
# Should show <100ms average

# Check concurrent sensor handling
go test -bench=BenchmarkConcurrentSensors ./tests/benchmark/
# Should handle 25 concurrent sensors
```

### 4. HTM Compliance Validation

Verify HTM biological constraints:

```bash
# Run HTM compliance tests
go test ./tests/contract -v -run TestHTMCompliance

# Check sparsity maintenance
go test ./tests/integration -v -run TestSparsityCompliance

# Validate temporal properties
go test ./tests/integration -v -run TestTemporalProperties

# Check biological constraints
go test ./tests/integration -v -run TestBiologicalConstraints
```

## Troubleshooting Common Issues

### 1. Spatial Pooler Integration Issues

If temporal memory fails to integrate with spatial pooler:

```bash
# Check spatial pooler output format
go test ./tests/contract -v -run TestSpatialPoolerOutput

# Verify SDR compatibility
go test ./tests/unit -v -run TestSDRCompatibility

# Check interface definitions
grep -r "SpatialPoolingService" internal/ports/
```

### 2. Motor Output Command Generation Issues

If motor output fails to generate commands:

```bash
# Check prediction format compatibility
go test ./tests/unit -v -run TestPredictionFormat

# Verify confidence thresholds
go test ./tests/contract -v -run TestMotorOutputThresholds

# Check action mapping configuration
go test ./tests/unit -v -run TestActionMappings
```

### 3. Sample Client API Communication Issues

If sample client cannot communicate with HTM API:

```bash
# Check API server status
curl -f http://localhost:8080/api/v1/health
# Should return 200 OK

# Verify endpoint availability
curl -f http://localhost:8080/api/v1/pipeline/status?pipeline_id=test
# Should return pipeline status

# Check client configuration
./sample_client --test-connection --verbose
# Should show detailed connection information
```

### 4. Performance Issues

If processing exceeds latency targets:

```bash
# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=BenchmarkCompletePipeline ./tests/benchmark/
go tool pprof cpu.prof

# Profile memory usage
go test -memprofile=mem.prof -bench=BenchmarkCompletePipeline ./tests/benchmark/
go tool pprof mem.prof

# Check for bottlenecks
go test -trace=trace.out -bench=BenchmarkCompletePipeline ./tests/benchmark/
go tool trace trace.out
```

## Next Steps

After completing the quickstart:

1. **Phase 2**: Run `/tasks` command to generate detailed implementation tasks
2. **Implementation**: Follow generated tasks to implement all components
3. **Validation**: Run comprehensive test suite to verify implementation
4. **Performance Tuning**: Optimize based on benchmark results
5. **Documentation**: Update README.md with complete pipeline examples

## Quick Reference

### Key Commands
- `go test ./tests/contract -v` - Run contract tests
- `go test ./tests/integration -v` - Run integration tests  
- `go test -bench=. ./tests/benchmark/` - Run performance benchmarks
- `go run cmd/api/main.go` - Start HTM API server
- `./sample_client --sensors=all --duration=60s` - Run complete test scenario

### Key Files
- `internal/cortical/temporal/` - Temporal memory implementation
- `internal/handlers/motor_output_handler.go` - Motor output HTTP handler
- `internal/handlers/pipeline_handler.go` - Complete pipeline handler
- `pkg/client/sample_client/` - Sample client application
- `tests/integration/complete_pipeline_test.go` - End-to-end tests

### Performance Targets
- End-to-end latency: <100ms
- Concurrent sensors: 25 simultaneous
- HTM sparsity: 2-5% maintained throughout pipeline
- Memory usage: <500MB for complete pipeline
- Success rate: >95% for known test scenarios

---

*Quickstart completed: October 2, 2025*  
*Ready for detailed task generation and implementation*