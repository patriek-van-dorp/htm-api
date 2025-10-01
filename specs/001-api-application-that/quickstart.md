# Quickstart Guide: HTM Neural Network API Core

**Feature**: 001-api-application-that  
**Date**: September 30, 2025  
**Estimated Completion Time**: 15 minutes

## Prerequisites

### Required Software
- Go 1.21 or later installed
- Git for version control
- curl or Postman for API testing
- (Optional) Docker for containerized deployment

### System Requirements
- Minimum 4GB RAM for development
- Network connectivity for dependency downloads
- Available port 8080 for local development

## Quick Setup

### 1. Initialize Go Module
```bash
# Create project directory
mkdir htm-neural-api
cd htm-neural-api

# Initialize Go module
go mod init github.com/your-org/htm-neural-api

# Create basic directory structure
mkdir -p cmd/api internal/{api/{handlers,middleware},domain/{htm,processing},infrastructure/{config,validation},ports} pkg/client tests/{contract,integration,unit}
```

### 2. Install Dependencies
```bash
# Core dependencies
go get github.com/gin-gonic/gin@latest
go get gonum.org/v1/gonum@latest
go get github.com/go-playground/validator/v10@latest
go get github.com/google/uuid@latest

# Testing dependencies
go get github.com/stretchr/testify@latest
go get github.com/golang/mock@latest
```

### 3. Create Main Application
Create `cmd/api/main.go`:
```go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize Gin router
    r := gin.Default()
    
    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })
    
    // Start server
    log.Println("Starting HTM Neural API on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

### 4. Run Basic Server
```bash
# Run the application
go run cmd/api/main.go

# Test health endpoint (in another terminal)
curl http://localhost:8080/health
```

**Expected Output**:
```json
{"status":"healthy"}
```

## API Usage Examples

### Example 1: Basic HTM Input Processing
```bash
# Send a simple 2D array for processing
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
    "request_id": "req-550e8400-e29b-41d4-a716-446655440001",
    "priority": "normal"
  }'
```

**Expected Response**:
```json
{
  "result": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "result": [[1.1, 2.1, 3.1], [4.1, 5.1, 6.1]],
    "metadata": {
      "processing_time_ms": 45,
      "instance_id": "api-001",
      "algorithm_version": "placeholder-v1.0"
    },
    "status": "SUCCESS"
  },
  "request_id": "req-550e8400-e29b-41d4-a716-446655440001",
  "response_time": "2025-09-30T10:00:00.045Z"
}
```

### Example 2: Invalid Input Handling
```bash
# Send invalid data (mismatched dimensions)
curl -X POST http://localhost:8080/api/v1/process \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "data": [[1.0, 2.0], [4.0, 5.0]],
      "metadata": {
        "dimensions": [2, 3],
        "sensor_id": "sensor001",
        "version": "v1.0"
      },
      "timestamp": "2025-09-30T10:00:00Z"
    },
    "request_id": "req-550e8400-e29b-41d4-a716-446655440002"
  }'
```

**Expected Response**:
```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "Data array dimensions do not match metadata",
    "details": {
      "expected_dimensions": [2, 3],
      "actual_dimensions": [2, 2]
    },
    "retryable": false
  },
  "request_id": "req-550e8400-e29b-41d4-a716-446655440002",
  "response_time": "2025-09-30T10:00:00.005Z"
}
```

### Example 3: Health and Metrics Monitoring
```bash
# Check service health
curl http://localhost:8080/health

# Get detailed metrics
curl http://localhost:8080/metrics
```

## Development Workflow

### 1. Test-Driven Development Setup
Create `tests/contract/api_test.go`:
```go
package contract

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestProcessHTMInput_Success(t *testing.T) {
    // This test should fail initially (no implementation)
    assert.Fail(t, "API endpoint not implemented yet")
}

func TestProcessHTMInput_InvalidInput(t *testing.T) {
    // This test should fail initially (no implementation)
    assert.Fail(t, "Input validation not implemented yet")
}
```

### 2. Run Tests (Should Fail Initially)
```bash
# Run contract tests
go test ./tests/contract/... -v

# Run all tests
go test ./... -v
```

### 3. Basic Project Validation
```bash
# Verify Go module
go mod tidy
go mod verify

# Check for common issues
go vet ./...

# Format code
go fmt ./...

# Build application
go build -o bin/htm-api cmd/api/main.go
```

## Performance Validation

### 1. Response Time Test
```bash
# Test response time requirement (<100ms acknowledgment)
time curl -X POST http://localhost:8080/api/v1/process \
  -H "Content-Type: application/json" \
  -d '{"input":{"id":"test","data":[[1]],"metadata":{"dimensions":[1,1],"sensor_id":"test","version":"v1.0"},"timestamp":"2025-09-30T10:00:00Z"},"request_id":"test"}'
```

### 2. Concurrent Request Test
```bash
# Simple concurrent test (run multiple times)
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/process \
    -H "Content-Type: application/json" \
    -d "{\"input\":{\"id\":\"test-$i\",\"data\":[[1]],\"metadata\":{\"dimensions\":[1,1],\"sensor_id\":\"test\",\"version\":\"v1.0\"},\"timestamp\":\"2025-09-30T10:00:00Z\"},\"request_id\":\"req-$i\"}" &
done
wait
```

## Container Deployment

### 1. Create Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/htm-api cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/htm-api .
EXPOSE 8080
CMD ["./htm-api"]
```

### 2. Build and Run Container
```bash
# Build Docker image
docker build -t htm-neural-api .

# Run container
docker run -p 8080:8080 htm-neural-api

# Test containerized API
curl http://localhost:8080/health
```

## Troubleshooting

### Common Issues

**Port Already in Use**:
```bash
# Find process using port 8080
lsof -i :8080

# Kill process if needed
kill -9 <PID>
```

**Module Not Found**:
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
```

**Build Failures**:
```bash
# Check Go version
go version

# Verify module integrity
go mod verify

# Update dependencies
go get -u ./...
```

### Development Tips
1. Use `go run cmd/api/main.go` for quick iteration
2. Enable Gin debug mode with `GIN_MODE=debug`
3. Use `go test -race` to detect race conditions
4. Monitor memory usage with `go tool pprof`

## Next Steps

After completing this quickstart:

1. **Implement Data Models**: Create Go structs in `internal/domain/`
2. **Add HTTP Handlers**: Implement API endpoints in `internal/api/handlers/`
3. **Matrix Processing**: Integrate gonum for mathematical operations
4. **Error Handling**: Implement robust error handling and retry logic
5. **Testing**: Expand test coverage for all components
6. **Monitoring**: Add metrics collection and health checks
7. **Documentation**: Generate API documentation from OpenAPI spec

## Success Criteria

✅ Basic server starts and responds to health checks  
✅ API endpoints accept JSON requests  
✅ Response format matches OpenAPI specification  
✅ Tests run (even if failing initially)  
✅ Container builds and runs successfully  
✅ Performance meets <100ms acknowledgment requirement  

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [gonum Documentation](https://pkg.go.dev/gonum.org/v1/gonum)
- [OpenAPI Specification](./contracts/openapi.yaml)
- [Feature Specification](./spec.md)
- [Data Model Documentation](./data-model.md)