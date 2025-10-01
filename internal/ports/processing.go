package ports

import (
	"context"

	"github.com/htm-project/neural-api/internal/domain/htm"
)

// ProcessingService defines the interface for HTM neural network processing operations.
type ProcessingService interface {
	// ProcessHTMInput processes HTM input data and returns the processing result
	ProcessHTMInput(ctx context.Context, input *htm.HTMInput) (*htm.ProcessingResult, error)

	// ValidateInput validates HTM input data without processing
	ValidateInput(ctx context.Context, input *htm.HTMInput) error

	// GetSupportedVersions returns the list of supported API versions
	GetSupportedVersions() []string

	// GetAlgorithmVersion returns the current algorithm version
	GetAlgorithmVersion() string

	// GetInstanceID returns the unique identifier for this processing instance
	GetInstanceID() string

	// HealthCheck performs a health check on the processing service
	HealthCheck(ctx context.Context) error
}

// MatrixProcessor defines the interface for matrix operations used in HTM processing.
type MatrixProcessor interface {
	// ProcessMatrix performs matrix processing operations on 2D data
	ProcessMatrix(ctx context.Context, data [][]float64) ([][]float64, error)

	// ValidateMatrix validates that a matrix is suitable for processing
	ValidateMatrix(data [][]float64) error

	// GetMatrixDimensions returns the dimensions of a matrix
	GetMatrixDimensions(data [][]float64) (rows, cols int)

	// IsMatrixConsistent checks if all rows have the same length
	IsMatrixConsistent(data [][]float64) bool

	// NormalizeMatrix normalizes matrix values to a standard range
	NormalizeMatrix(data [][]float64) ([][]float64, error)

	// CreateEmptyMatrix creates a matrix with specified dimensions
	CreateEmptyMatrix(rows, cols int) [][]float64
}

// ValidationService defines the interface for input validation operations.
type ValidationService interface {
	// ValidateHTMInput validates complete HTM input structure
	ValidateHTMInput(input *htm.HTMInput) error

	// ValidateAPIRequest validates complete API request structure
	ValidateAPIRequest(request *htm.APIRequest) error

	// ValidateMetadata validates input metadata
	ValidateMetadata(metadata *htm.InputMetadata) error

	// ValidateDimensions validates that declared dimensions match actual data
	ValidateDimensions(data [][]float64, dimensions []int) error

	// ValidateUUID validates UUID format
	ValidateUUID(uuid string) error

	// ValidateSensorID validates sensor ID format
	ValidateSensorID(sensorID string) error
}

// MetricsCollector defines the interface for collecting processing metrics.
type MetricsCollector interface {
	// IncrementRequestCount increments the total request counter
	IncrementRequestCount()

	// IncrementErrorCount increments the error counter
	IncrementErrorCount()

	// RecordProcessingTime records the time taken for processing
	RecordProcessingTime(duration int64)

	// RecordResponseTime records the total response time
	RecordResponseTime(duration int64)

	// SetConcurrentRequests sets the current number of concurrent requests
	SetConcurrentRequests(count int)

	// GetMetrics returns current metrics snapshot
	GetMetrics() map[string]interface{}

	// Reset resets all metrics
	Reset()
}

// ProcessingRepository defines the interface for persisting processing results (future use).
type ProcessingRepository interface {
	// SaveResult saves a processing result (for audit/debugging)
	SaveResult(ctx context.Context, result *htm.ProcessingResult) error

	// GetResult retrieves a processing result by ID
	GetResult(ctx context.Context, id string) (*htm.ProcessingResult, error)

	// GetResultsByTimeRange retrieves results within a time range
	GetResultsByTimeRange(ctx context.Context, start, end int64) ([]*htm.ProcessingResult, error)

	// DeleteOldResults removes results older than specified timestamp
	DeleteOldResults(ctx context.Context, olderThan int64) error
}
