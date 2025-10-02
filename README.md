# HTM Neural Network API - Cortical Processing Pipeline

A high-performance REST API for processing HTM (Hierarchical Temporal Memory) neural networks, implementing the complete cortical column architecture from sensor input to motor output. Built with Go for optimal concurrent request handling and matrix operations using gonum.

## HTM Architecture Overview

This API implements the biological cortical processing pipeline:

```
Sensor Data → Encoders → Spatial Pooler → Temporal Memory → Motor Output
     ↓            ↓           ↓              ↓             ↓
 Raw Values → Bit Patterns → HTM SDRs → Sequences → Actions
```

## Features

- **Complete HTM Pipeline**: Full cortical column implementation from input to output
- **Biological Constraints**: Enforces proper sparsity (2-5%) and temporal patterns
- **Spatial Pooler**: Converts encoder outputs to HTM-compliant sparse representations
- **Temporal Memory**: Learns sequences and makes predictions (future enhancement)
- **Motor Output**: Converts predictions to actionable outputs (future enhancement)
- **Fast Processing**: <50ms per cortical processing stage
- **Concurrent Support**: Handle multiple simultaneous cortical processing requests
- **Real-time Monitoring**: Health checks, metrics, and HTM properties validation
- **Matrix Operations**: Efficient multi-dimensional array processing with gonum

## Quick Start

### Prerequisites

- Go 1.23+ installed
- Git for version control

### Installation

```bash
# Clone repository
git clone https://github.com/htm-project/neural-api.git
cd neural-api

# Install dependencies
go mod download

# Run the API server
go run cmd/api/main.go
```

The server will start on `http://localhost:8080` with the following message:
```
Starting HTM Neural Processing API server on 0.0.0.0:8080
Environment: development
Debug mode: true
```

## API Usage Guide

### Complete HTM Processing Pipeline

The API provides access to different stages of cortical processing. In a complete HTM system, data flows through multiple stages:

#### Stage 1: Sensor Input → Encoder (Future Enhancement)
```bash
# Future: Process raw sensor data through encoders
curl -X POST http://localhost:8080/api/v1/encode \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_data": {
      "type": "temperature",
      "value": 72.5,
      "units": "fahrenheit"
    },
    "encoder_config": {
      "type": "numeric",
      "min": 32.0,
      "max": 100.0,
      "resolution": 0.1
    }
  }'
```

#### Stage 2: Spatial Pooler Processing (Current Implementation)
**Purpose**: Converts encoder outputs to HTM-compliant SDRs for temporal memory processing

```bash
curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "encoder_output": {
      "width": 1024,
      "active_bits": [0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210],
      "sparsity": 0.021
    },
    "input_width": 1024,
    "input_id": "sensor-reading-001",
    "learning_enabled": true,
    "metadata": {
      "sensor_type": "temperature",
      "location": "room_a"
    }
  }'
```

**Response** (HTM-compliant SDR for temporal memory):
```json
{
  "normalized_sdr": {
    "width": 2048,
    "active_bits": [5, 12, 23, 45, 67, 89, 123, 156, 178, 234, 267, 345, 456, 567, 678, 789, 890, 901, 1012, 1123, 1234, 1345, 1456, 1567, 1678, 1789, 1890, 1901, 2012],
    "sparsity": 0.025
  },
  "input_id": "sensor-reading-001",
  "processing_time_ms": 15,
  "active_columns": 51,
  "avg_overlap": 8.5,
  "sparsity_level": 0.025,
  "learning_occurred": true,
  "boosting_applied": false
}
```

#### Stage 3: Temporal Memory Processing (Future Enhancement)
```bash
# Future: Process spatial pooler SDRs through temporal memory
curl -X POST http://localhost:8080/api/v1/temporal-memory/process \
  -H "Content-Type: application/json" \
  -d '{
    "spatial_sdr": {
      "width": 2048,
      "active_bits": [5, 12, 23, 45, 67, 89, 123, 156, 178, 234, 267, 345, 456, 567, 678, 789, 890, 901, 1012, 1123, 1234, 1345, 1456, 1567, 1678, 1789, 1890, 1901, 2012],
      "sparsity": 0.025
    },
    "learning_enabled": true
  }'
```

#### Stage 4: Motor Output Processing (Future Enhancement)
```bash
# Future: Convert predictions to motor actions
curl -X POST http://localhost:8080/api/v1/motor-output/process \
  -H "Content-Type: application/json" \
  -d '{
    "prediction_sdr": {
      "width": 2048,
      "active_bits": [...],
      "confidence": 0.85
    },
    "motor_config": {
      "output_type": "movement",
      "action_space": ["forward", "backward", "left", "right", "stop"]
    }
  }'
```

### Current Implementation: Spatial Pooler Stage

### 1. Health Check

Check if the cortical processing pipeline is running correctly:

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "dependencies": {},
  "service": {
    "processing_service": true,
    "spatial_pooling_service": true,
    "metrics_collector": true,
    "spatial_pooler_healthy": true,
    "uptime_seconds": 1.234
  },
  "system": {
    "memory_mb": 45.2,
    "goroutines": 8,
    "gc_pause_ms": 0.5
  },
  "healthy": true
}
```

### 2. Spatial Pooler Status

Get detailed status of the HTM spatial pooler (cortical processing stage):

```bash
curl http://localhost:8080/api/v1/spatial-pooler/status
```

**Response:**
```json
{
  "status": "operational",
  "healthy": true,
  "instance": {
    "instance_id": "main-instance",
    "created_at": "2025-10-02T10:00:00Z",
    "uptime_seconds": 123.45
  },
  "configuration": {
    "input_width": 1024,
    "column_count": 2048,
    "sparsity_ratio": 0.025,
    "learning_enabled": true,
    "mode": "deterministic"
  },
  "metrics": {
    "total_processed": 0,
    "average_processing_time_ms": 0,
    "learning_iterations": 0,
    "average_sparsity": 0
  },
  "timestamp": "2025-10-02T10:00:00Z"
}
```

### 3. HTM Pipeline Integration

**Important**: The spatial pooler output is designed for the next stage in HTM processing (temporal memory), not for direct sensor feedback. The API provides access to intermediate stages for debugging, research, and modular HTM system development.

```bash
curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "encoder_output": {
      "width": 1024,
      "active_bits": [0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210],
      "sparsity": 0.021
    },
    "input_width": 1024,
    "input_id": "sensor-reading-001",
    "learning_enabled": true,
    "metadata": {
      "sensor_type": "temperature",
      "location": "room_a"
    }
  }'
```

**Response:**
```json
{
  "normalized_sdr": {
    "width": 2048,
    "active_bits": [5, 12, 23, 45, 67, 89, 123, 156, 178, 234, 267, 345, 456, 567, 678, 789, 890, 901, 1012, 1123, 1234, 1345, 1456, 1567, 1678, 1789, 1890, 1901, 2012],
    "sparsity": 0.025
  },
  "input_id": "sensor-reading-001",
  "processing_time_ms": 15,
  "active_columns": 51,
  "avg_overlap": 8.5,
  "sparsity_level": 0.025,
  "learning_occurred": true,
  "boosting_applied": false
}
```

### 4. HTM Properties Validation

Validate HTM biological properties and cortical compliance:

```bash
curl http://localhost:8080/api/v1/spatial-pooler/validation/htm-properties
```

**Response:**
```json
{
  "htm_compliance": {
    "biological_constraints": {
      "sparsity_percentage": 2.5,
      "target_sparsity_range": [2.0, 5.0],
      "sparsity_compliant": true,
      "overlap_threshold": 5,
      "overlap_compliant": true
    },
    "learning_properties": {
      "learning_enabled": true,
      "learning_rate": 0.1,
      "boost_strength": 0.5,
      "learning_compliant": true
    },
    "topology_properties": {
      "column_count": 2048,
      "input_width": 1024,
      "inhibition_radius": 16,
      "topology_compliant": true
    }
  },
  "runtime_metrics": {
    "current_sparsity": 2.5,
    "total_processed": 5,
    "average_processing_time": 18,
    "learning_iterations": 5
  },
  "validation_status": {
    "overall_compliant": true,
    "warnings": [],
    "recommendations": []
  }
}
```

### 5. Configuration Management

Get current spatial pooler configuration:

```bash
curl http://localhost:8080/api/v1/spatial-pooler/config
```

Update spatial pooler configuration:

```bash
curl -X PUT http://localhost:8080/api/v1/spatial-pooler/config \
  -H "Content-Type: application/json" \
  -d '{
    "sparsity_ratio": 0.03,
    "learning_rate": 0.15,
    "boost_strength": 0.6
  }'
```

### 6. Metrics and Monitoring

Get spatial pooler processing metrics:

```bash
curl http://localhost:8080/api/v1/spatial-pooler/metrics
```

Reset metrics:

```bash
curl -X POST http://localhost:8080/api/v1/spatial-pooler/metrics/reset
```

## HTM Cortical Processing Architecture

### Complete HTM System Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Sensors   │    │  Encoders   │    │  Spatial    │    │  Temporal   │    │   Motor     │
│             │───▶│             │───▶│   Pooler    │───▶│   Memory    │───▶│   Output    │
│ (Raw Data)  │    │(Bit Pattern)│    │ (HTM SDRs)  │    │(Sequences)  │    │ (Actions)   │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
      │                    │                    │                    │                    │
   Examples:           Converts to:        Learns spatial       Learns temporal     Generates:
 - Temperature        - Dense/sparse       patterns &           sequences &         - Movement
 - Images             bit arrays          enforces HTM         makes predictions   - Sound  
 - Text               - 1024-2048 bits    biological           - Anticipation      - Control
 - Categories         - Variable sparsity constraints          - Context           signals
                                         - 2-5% sparsity      - Anomaly detection
```

### Current Implementation Status

| Component | Status | Purpose |
|-----------|--------|---------|
| **Sensors** | 🔄 Planned | Raw data input (temperature, images, text, etc.) |
| **Encoders** | 🔄 Planned | Convert raw data to bit patterns |
| **Spatial Pooler** | ✅ Implemented | Convert bit patterns to HTM-compliant SDRs |
| **Temporal Memory** | 🔄 Planned | Learn sequences and make predictions |
| **Motor Output** | 🔄 Planned | Convert predictions to actions |

### Why This Architecture Matters

1. **Spatial Pooler Role**: 
   - **Input**: Dense/sparse bit patterns from encoders
   - **Output**: HTM-compliant SDRs (2-5% sparsity)
   - **Purpose**: Normalize representations for temporal processing
   - **Next Stage**: Feeds temporal memory, NOT back to sensors

2. **Temporal Memory Role** (future):
   - **Input**: HTM SDRs from spatial pooler
   - **Output**: Sequence predictions and anomaly detection
   - **Purpose**: Learn temporal patterns and make predictions
   - **Next Stage**: Feeds motor output systems

3. **Motor Output Role** (future):
   - **Input**: Predictions from temporal memory
   - **Output**: Physical actions, control signals, responses
   - **Purpose**: Convert neural predictions to real-world actions
   - **Examples**: Robot movement, speech generation, system control

## Testing Guide

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test suites
go test ./tests/contract/        # API contract tests
go test ./tests/integration/     # End-to-end integration tests
go test ./tests/unit/           # Unit tests
```

### Contract Testing

Test specific spatial pooler functionality:

```bash
# Test spatial pooler status endpoint
go test -v -run TestSpatialPoolerStatusBasic ./tests/contract/spatial_pooler_status_simple_test.go

# Test spatial pooler processing
go test -v -run TestSpatialPoolerProcessBasic ./tests/contract/spatial_pooler_process_simple_test.go

# Test HTM validation
go test -v -run TestHTMValidation ./tests/contract/htm_validation_test.go
```

### Integration Testing

Test complete pipelines:

```bash
# Test complete HTM processing pipeline
go test -v -run TestCompletePipeline ./tests/integration/complete_pipeline_test.go

# Test performance requirements
go test -v -run TestPerformance ./tests/integration/performance_integration_test.go

# Test concurrent processing
go test -v -run TestConcurrency ./tests/integration/concurrency_integration_test.go
```

### Manual Testing with curl

#### Test Valid Input (Should Succeed)

```bash
curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "encoder_output": {
      "width": 1024,
      "active_bits": [0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210],
      "sparsity": 0.025
    },
    "input_width": 1024,
    "input_id": "test-valid-input",
    "learning_enabled": true
  }'
```

#### Test Invalid Sparsity (Should Fail)

```bash
curl -X POST http://localhost:8080/api/v1/spatial-pooler/process \
  -H "Content-Type: application/json" \
  -d '{
    "encoder_output": {
      "width": 1024,
      "active_bits": [0, 10],
      "sparsity": 0.001
    },
    "input_width": 1024,
    "input_id": "test-invalid-sparsity",
    "learning_enabled": true
  }'
```

**Expected Error Response:**
```json
{
  "error": "Spatial pooling processing failed",
  "details": "spatial pooling failed: processing_error: sparsity level 0.001 is outside HTM range [0.02, 0.05]"
}
```

## HTM Spatial Pooler Configuration

**Important Note**: The spatial pooler is configured for **internal cortical processing**, not end-user applications. Its output feeds the temporal memory stage in a complete HTM system.

### Key Parameters

- **Sparsity Ratio** (0.02-0.05): Percentage of active columns (HTM biological constraint)
- **Learning Rate** (0.0-1.0): Rate of synaptic adaptation
- **Column Count**: Number of cortical columns (default: 2048)
- **Input Width**: Size of input space (default: 1024)
- **Inhibition Radius**: Local inhibition neighborhood size
- **Boost Strength**: Column usage balancing factor

### Example Configuration

```json
{
  "input_width": 1024,
  "column_count": 2048,
  "sparsity_ratio": 0.025,
  "mode": "deterministic",
  "learning_enabled": true,
  "learning_rate": 0.1,
  "max_boost": 3.0,
  "boost_strength": 0.5,
  "inhibition_radius": 16,
  "local_area_density": 0.025,
  "min_overlap_threshold": 5,
  "max_processing_time_ms": 50
```

## Development

### Project Structure

```
cmd/
  api/                     # Application entry point
    main.go               # Server initialization with dependency injection
internal/
  api/
    router.go            # HTTP routing and middleware setup
  cortical/
    spatial/             # HTM spatial pooler implementation
  handlers/
    spatial_pooler_handler.go    # Spatial pooler HTTP endpoints
    health_metrics_handler.go    # System health and metrics
  infrastructure/
    config/              # Configuration management
  services/
    spatial_pooling_service.go   # HTM business logic
    processing_service.go        # General processing
tests/
  contract/              # API contract validation tests
  integration/           # End-to-end feature tests
  unit/                 # Component unit tests
```

### Development Commands

```bash
# Install development dependencies
go mod download
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Code quality checks
golangci-lint run
go vet ./...
go fmt ./...

# Run tests with verbose output
go test -v ./...

# Run specific test pattern
go test -v -run "TestSpatial" ./...

# Build for production
go build -o htm-api cmd/api/main.go

# Run with custom configuration
export HTM_PORT=9090
export HTM_DEBUG=false
go run cmd/api/main.go
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HTM_PORT` | `8080` | Server port |
| `HTM_HOST` | `0.0.0.0` | Server bind address |
| `HTM_DEBUG` | `true` | Enable debug logging |
| `HTM_SPARSITY` | `0.025` | Spatial pooler sparsity ratio |
| `HTM_COLUMNS` | `2048` | Number of cortical columns |
| `HTM_INPUT_WIDTH` | `1024` | Input space width |
| `HTM_LEARNING_RATE` | `0.1` | Synaptic learning rate |

### Adding New Features

1. **New Endpoints**: Add handlers in `internal/handlers/`
2. **Business Logic**: Implement services in `internal/services/`
3. **Configuration**: Update `internal/infrastructure/config/`
4. **Tests**: Add corresponding tests in `tests/`

## Troubleshooting

### Common Issues and Solutions

#### 1. Sparsity Constraint Violations

**Error:**
```
spatial pooling failed: processing_error: sparsity level 0.001 is outside HTM range [0.02, 0.05]
```

**Solution:**
Ensure input sparsity is between 2-5%. Adjust your encoder to produce more active bits:

```bash
# Check current input sparsity
total_bits = len(active_bits)
required_bits = int(input_width * 0.025)  # For 2.5% sparsity
```

#### 2. Processing Time Violations

**Error:**
```
spatial pooling failed: processing_error: processing time 75ms exceeds maximum allowed 50ms
```

**Solution:**
- Reduce input width or column count
- Check system resources (CPU/memory)
- Increase timeout in configuration

#### 3. Server Won't Start

**Error:**
```
Failed to start server: listen tcp :8080: bind: address already in use
```

**Solution:**
```bash
# Check what's using port 8080
netstat -ano | findstr :8080

# Use different port
export HTM_PORT=8081
go run cmd/api/main.go
```

#### 4. Memory Issues with Large Inputs

**Error:**
```
runtime: out of memory
```

**Solution:**
- Reduce matrix dimensions
- Process in smaller batches
- Check input width and column count ratio

### Debug Mode

Enable detailed logging:

```bash
export HTM_DEBUG=true
go run cmd/api/main.go
```

Debug output includes:
- Request/response details
- HTM processing metrics
- Matrix operation timings
- Memory usage statistics

### Performance Tuning

#### Optimal Configuration for Different Use Cases

**High Throughput (Low Latency):**
```json
{
  "input_width": 512,
  "column_count": 1024,
  "sparsity_ratio": 0.02,
  "max_processing_time_ms": 25
}
```

**High Accuracy (Learning Focus):**
```json
{
  "input_width": 2048,
  "column_count": 4096,
  "sparsity_ratio": 0.03,
  "learning_rate": 0.05,
  "max_processing_time_ms": 100
}
```

**Balanced (General Purpose):**
```json
{
  "input_width": 1024,
  "column_count": 2048,
  "sparsity_ratio": 0.025,
  "learning_rate": 0.1,
  "max_processing_time_ms": 50
}
```

## HTM Theory Background

### Hierarchical Temporal Memory (HTM)

HTM is a biologically-inspired machine learning algorithm that models the structure and function of the neocortex. **The spatial pooler is the first cortical processing stage**, not a standalone application.

#### Complete HTM Pipeline

```
Input Data → Encoding → Spatial Pooling → Temporal Memory → Motor Output
     │           │            │               │              │
  Raw values → Bit patterns → HTM SDRs → Sequences → Actions
```

#### Key Principles

1. **Sparse Distributed Representations (SDRs)**
   - Only 2-5% of neurons are active at any time
   - High dimensional, sparse binary vectors
   - Robust to noise and fault-tolerant

2. **Spatial Pooling** (Current Implementation)
   - Converts encoder outputs to sparse representations
   - Learns spatial patterns through competitive learning
   - **Maintains semantic similarity through overlap**
   - **Prepares data for temporal processing**

3. **Temporal Memory** (Future Implementation)
   - Learns sequences of spatial patterns
   - Makes predictions about future inputs
   - Detects anomalies and novelty
   - **Feeds predictions to motor output**

4. **Motor Output** (Future Implementation)
   - Converts neural predictions to actions
   - Controls physical systems, robots, interfaces
   - Provides feedback to the learning system

#### Biological Constraints

- **Sparsity levels**: Mimic cortical column activation (2-5%)
- **Local inhibition**: Simulates lateral inhibition in cortex
- **Synaptic plasticity**: Follows Hebbian learning principles
- **Sequential processing**: Temporal patterns emerge naturally

#### Spatial Pooler Algorithm (Current Stage)

```
For each encoder input:
1. Compute overlap between input and column synapses
2. Apply boosting to underrepresented columns  
3. Perform local inhibition to maintain sparsity
4. Update synaptic permanences (learning)
5. Output HTM-compliant SDR for temporal memory
```

#### Benefits of Complete HTM System

- **Noise Resistance**: Sparse representations are robust to input corruption
- **Semantic Preservation**: Similar inputs produce similar SDRs with measurable overlap
- **Continuous Learning**: Online adaptation without catastrophic forgetting
- **Biological Plausibility**: Based on neuroscience principles and cortical architecture
- **Predictive Capability**: Temporal memory enables anticipation and anomaly detection
- **Action Generation**: Motor output converts predictions to real-world responses

### HTM Processing Stages

| Stage | Input | Output | Current Status | Biological Analog |
|-------|-------|--------|----------------|-------------------|
| **Encoding** | Raw sensor data | Bit patterns | 🔄 Planned | Sensory receptors |
| **Spatial Pooler** | Bit patterns | HTM SDRs | ✅ Implemented | Layer 4 cortical columns |
| **Temporal Memory** | HTM SDRs | Predictions | 🔄 Planned | Layer 2/3 pyramidal cells |
| **Motor Output** | Predictions | Actions | 🔄 Planned | Motor cortex output |

### Spatial Pooler Parameters

| Parameter | Range | Description | HTM Principle |
|-----------|-------|-------------|---------------|
| Sparsity | 2-5% | Active column percentage | Cortical sparsity |
| Learning Rate | 0.0-1.0 | Synaptic adaptation speed | Hebbian plasticity |
| Inhibition Radius | 5-50 | Local competition area | Lateral inhibition |
| Boost Strength | 0.0-3.0 | Usage balancing factor | Homeostatic regulation |

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make changes and add tests
4. Run test suite: `go test ./...`
5. Commit changes: `git commit -m 'Add amazing feature'`
6. Push to branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Code Standards

- Follow Go conventions and idioms
- Add tests for new functionality
- Update documentation for API changes
- Ensure HTM biological constraints are maintained
- Include performance benchmarks for spatial pooler changes

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## References

- [HTM School](https://numenta.com/htm-school/) - HTM theory and algorithms
- [Numenta Research](https://numenta.com/papers-videos-and-more/) - HTM papers and publications
- [Spatial Pooler Algorithm](https://numenta.com/spatial-pooler/) - Detailed spatial pooler explanation
- [Go Documentation](https://golang.org/doc/) - Go programming language
- [Gonum](https://www.gonum.org/) - Go numerical computing libraries

# Build with optimizations
go build -ldflags="-w -s" -o bin/htm-api cmd/api/main.go
```

### Docker

```bash
# Build image
docker build -t htm-neural-api .

# Run container
docker run -p 8080:8080 htm-neural-api
```

## Architecture

The API follows hexagonal architecture principles:

- **Domain Layer**: Core business entities and processing logic
- **Ports**: Interfaces defining contracts between layers
- **Adapters**: Implementations of ports (HTTP handlers, services)
- **Infrastructure**: External concerns (configuration, validation)

## Performance Goals

- **Response Time**: <100ms acknowledgment for processing requests
- **Concurrency**: Support multiple simultaneous requests from different sensors
- **Scalability**: Horizontal scaling through stateless design
- **Memory**: Efficient matrix operations with gonum library

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Related Documentation

- [Implementation Plan](specs/001-api-application-that/plan.md)
- [Data Model](specs/001-api-application-that/data-model.md)
- [API Contracts](specs/001-api-application-that/contracts/)
- [Research Notes](specs/001-api-application-that/research.md)