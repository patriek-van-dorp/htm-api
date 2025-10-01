# HTM Neural Network API Core

A high-performance REST API for processing HTM (Hierarchical Temporal Memory) neural network inputs. Built with Go for optimal concurrent request handling and matrix operations.

## Features

- **Fast Processing**: <100ms acknowledgment for API requests
- **Concurrent Support**: Handle multiple sensor inputs simultaneously  
- **API Chaining**: Consistent input/output format for pipeline integration
- **Scalable Architecture**: Stateless design for horizontal scaling
- **Matrix Operations**: Efficient multi-dimensional array processing with gonum

## Quick Start

### Prerequisites

- Go 1.21+ installed
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

### Test the API

```bash
# Health check
curl http://localhost:8080/health

# Process HTM input
curl -X POST http://localhost:8080/api/v1/process \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "data": [[1.0, 2.0, 3.0], [4.0, 5.0, 6.0]],
      "metadata": {
        "dimensions": [2, 3],
        "sensor_id": "sensor001",
        "version": "v1.0"
      },
      "timestamp": "2025-09-30T10:00:00Z"
    },
    "request_id": "req-550e8400-e29b-41d4-a716-446655440001"
  }'
```

## Project Structure

```
├── cmd/api/                    # Application entry point
├── internal/
│   ├── api/                   # HTTP layer (handlers, middleware, routing)
│   ├── domain/                # Business logic and entities
│   ├── infrastructure/        # External concerns (config, validation)
│   └── ports/                 # Interfaces for dependency inversion
├── pkg/client/                # Go client library
└── tests/                     # Test suites
    ├── contract/              # API contract tests
    ├── integration/           # End-to-end integration tests
    └── unit/                  # Unit tests
```

## API Documentation

### Endpoints

- `POST /api/v1/process` - Process HTM neural network input
- `GET /health` - Health check endpoint  
- `GET /metrics` - Performance metrics

See `specs/001-api-application-that/contracts/openapi.yaml` for full API specification.

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test suite
go test ./tests/contract/
go test ./tests/integration/
go test ./tests/unit/
```

### Building

```bash
# Build binary
go build -o bin/htm-api cmd/api/main.go

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