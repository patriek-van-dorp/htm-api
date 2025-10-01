# Data Model: HTM Neural Network API Core

**Feature**: 001-api-application-that  
**Date**: September 30, 2025  
**Status**: Complete

## Entity Definitions

### HTMInput
Represents multi-dimensional arrays containing spatial-temporal pattern data for HTM neural network processing.

**Fields**:
- `Data` ([][]float64): Multi-dimensional array representing spatial-temporal patterns
- `Metadata` (InputMetadata): Associated metadata for processing context
- `Timestamp` (time.Time): When the input was created/received
- `ID` (string): Unique identifier for tracking and correlation

**Validation Rules**:
- Data array must not be empty
- All sub-arrays must have consistent dimensions
- Metadata must include valid dimensions
- Timestamp must be valid RFC3339 format
- ID must be non-empty UUID format

**Go Type Definition**:
```go
type HTMInput struct {
    ID        string          `json:"id" validate:"required,uuid"`
    Data      [][]float64     `json:"data" validate:"required,min=1"`
    Metadata  InputMetadata   `json:"metadata" validate:"required"`
    Timestamp time.Time       `json:"timestamp" validate:"required"`
}
```

### InputMetadata
Provides context and configuration for HTM input processing.

**Fields**:
- `Dimensions` ([]int): Shape of the data array [rows, cols, ...]
- `SensorID` (string): Identifier for the source sensor/system
- `ProcessingHints` (map[string]interface{}): Optional processing parameters
- `Version` (string): API version for compatibility

**Validation Rules**:
- Dimensions must match actual data array shape
- SensorID must be non-empty alphanumeric string
- Version must match supported API versions (v1.0)
- ProcessingHints are optional but must be valid JSON if present

**Go Type Definition**:
```go
type InputMetadata struct {
    Dimensions      []int                  `json:"dimensions" validate:"required,min=2"`
    SensorID        string                 `json:"sensor_id" validate:"required,alphanum"`
    ProcessingHints map[string]interface{} `json:"processing_hints,omitempty"`
    Version         string                 `json:"version" validate:"required,oneof=v1.0"`
}
```

### ProcessingResult
Represents the output from HTM neural network computation, maintaining same format as input for API chaining.

**Fields**:
- `Result` ([][]float64): Processed multi-dimensional array (same format as input)
- `Metadata` (ResultMetadata): Processing metadata and performance info
- `Status` (ProcessingStatus): Success/failure status of processing
- `ID` (string): Correlation ID matching input ID

**Validation Rules**:
- Result array must maintain same dimensions as input
- Metadata must include processing time and instance ID
- Status must be valid enum value
- ID must match corresponding input ID

**Go Type Definition**:
```go
type ProcessingResult struct {
    ID       string           `json:"id" validate:"required,uuid"`
    Result   [][]float64      `json:"result" validate:"required"`
    Metadata ResultMetadata   `json:"metadata" validate:"required"`
    Status   ProcessingStatus `json:"status" validate:"required"`
}
```

### ResultMetadata
Contains processing performance and context information.

**Fields**:
- `ProcessingTimeMs` (int64): Time taken for processing in milliseconds
- `InstanceID` (string): API instance that processed the request
- `AlgorithmVersion` (string): Version of HTM algorithm used (placeholder for future)
- `QualityMetrics` (map[string]float64): Optional quality/confidence metrics

**Validation Rules**:
- ProcessingTimeMs must be non-negative
- InstanceID must be non-empty
- AlgorithmVersion defaults to "placeholder-v1.0"
- QualityMetrics values must be between 0.0 and 1.0

**Go Type Definition**:
```go
type ResultMetadata struct {
    ProcessingTimeMs   int64              `json:"processing_time_ms" validate:"min=0"`
    InstanceID         string             `json:"instance_id" validate:"required"`
    AlgorithmVersion   string             `json:"algorithm_version" validate:"required"`
    QualityMetrics     map[string]float64 `json:"quality_metrics,omitempty"`
}
```

### ProcessingStatus
Enumeration for processing result status.

**Values**:
- `SUCCESS`: Processing completed successfully
- `PARTIAL_SUCCESS`: Processing completed with warnings
- `FAILED`: Processing failed
- `TIMEOUT`: Processing exceeded time limits
- `RETRYING`: Processing is being retried

**Go Type Definition**:
```go
type ProcessingStatus string

const (
    StatusSuccess        ProcessingStatus = "SUCCESS"
    StatusPartialSuccess ProcessingStatus = "PARTIAL_SUCCESS"
    StatusFailed         ProcessingStatus = "FAILED"
    StatusTimeout        ProcessingStatus = "TIMEOUT"
    StatusRetrying       ProcessingStatus = "RETRYING"
)
```

### APIRequest
Wrapper for incoming HTTP requests with validation and context.

**Fields**:
- `Input` (HTMInput): The HTM input data to process
- `RequestID` (string): Unique request identifier for tracing
- `ClientID` (string): Optional client identifier
- `Priority` (RequestPriority): Processing priority level

**Validation Rules**:
- Input must pass all HTMInput validation rules
- RequestID must be unique UUID
- ClientID is optional but must be alphanumeric if provided
- Priority must be valid enum value

**Go Type Definition**:
```go
type APIRequest struct {
    Input     HTMInput        `json:"input" validate:"required"`
    RequestID string          `json:"request_id" validate:"required,uuid"`
    ClientID  string          `json:"client_id,omitempty" validate:"omitempty,alphanum"`
    Priority  RequestPriority `json:"priority" validate:"omitempty,oneof=low normal high"`
}
```

### APIResponse
Wrapper for HTTP responses with error handling.

**Fields**:
- `Result` (*ProcessingResult): Processing result (nil on error)
- `Error` (*APIError): Error information (nil on success)
- `RequestID` (string): Correlation ID from request
- `ResponseTime` (time.Time): When response was generated

**Validation Rules**:
- Either Result or Error must be non-nil, but not both
- RequestID must match incoming request
- ResponseTime must be valid timestamp

**Go Type Definition**:
```go
type APIResponse struct {
    Result       *ProcessingResult `json:"result,omitempty"`
    Error        *APIError         `json:"error,omitempty"`
    RequestID    string            `json:"request_id" validate:"required"`
    ResponseTime time.Time         `json:"response_time" validate:"required"`
}
```

### APIError
Structured error information for API responses.

**Fields**:
- `Code` (string): Error code for programmatic handling
- `Message` (string): Human-readable error description
- `Details` (map[string]interface{}): Additional error context
- `Retryable` (bool): Whether the request can be retried

**Go Type Definition**:
```go
type APIError struct {
    Code      string                 `json:"code" validate:"required"`
    Message   string                 `json:"message" validate:"required"`
    Details   map[string]interface{} `json:"details,omitempty"`
    Retryable bool                   `json:"retryable"`
}
```

## Relationships

### Processing Flow
1. `APIRequest` contains `HTMInput` with `InputMetadata`
2. Processing generates `ProcessingResult` with `ResultMetadata`
3. `APIResponse` wraps either `ProcessingResult` or `APIError`
4. All entities linked by correlation IDs for tracing

### Data Flow
```
HTMInput → Matrix Processing → ProcessingResult
    ↓              ↓                    ↓
InputMetadata → Processing → ResultMetadata
    ↓              ↓                    ↓
APIRequest → HTTP Handler → APIResponse
```

## State Transitions

### ProcessingStatus Flow
```
[Initial] → SUCCESS (happy path)
[Initial] → FAILED (validation/processing error)
[Initial] → TIMEOUT → RETRYING → SUCCESS/FAILED
[Initial] → RETRYING → SUCCESS/FAILED (transient error)
```

### Request Lifecycle
```
APIRequest → Validation → Processing → APIResponse
     ↓           ↓            ↓           ↓
  Received → Validated → Processing → Complete
                ↓            ↓
             Invalid → Error Response
                         ↓
                    Failed → Retry (if retryable)
```

## Validation Strategy

### Input Validation
- JSON schema validation for structure
- Go struct tags for field validation
- Custom validators for business rules (dimension consistency)
- Range validation for numeric values

### Output Validation
- Dimension consistency between input and output
- Response time SLA validation (<100ms acknowledgment)
- Error response format validation

### Performance Constraints
- Maximum input array size limits
- Processing timeout thresholds
- Memory usage monitoring per request

## Future Extensions

### HTM Algorithm Integration
Current data model designed to accommodate future HTM-specific fields:
- Sparse representation metadata
- Temporal sequence context
- Hierarchical level information
- Learning state parameters

### Enhanced Metadata
Planned extensions for advanced features:
- Confidence scores
- Processing quality metrics
- Algorithm-specific parameters
- Performance profiling data