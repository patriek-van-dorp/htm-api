package services

import (
	"context"
	"fmt"
	"time"

	"github.com/htm-project/neural-api/internal/domain/htm"
	"github.com/htm-project/neural-api/internal/ports"
)

// ProcessingServiceImpl implements the ProcessingService interface.
type ProcessingServiceImpl struct {
	matrixProcessor   ports.MatrixProcessor
	validationService ports.ValidationService
	metricsCollector  ports.MetricsCollector
	instanceID        string
}

// NewProcessingService creates a new processing service.
func NewProcessingService(
	matrixProcessor ports.MatrixProcessor,
	validationService ports.ValidationService,
	metricsCollector ports.MetricsCollector,
) ports.ProcessingService {
	return &ProcessingServiceImpl{
		matrixProcessor:   matrixProcessor,
		validationService: validationService,
		metricsCollector:  metricsCollector,
		instanceID:        "processing-service-1",
	}
}

// ProcessHTMInput processes HTM input data and returns the processing result.
func (ps *ProcessingServiceImpl) ProcessHTMInput(ctx context.Context, input *htm.HTMInput) (*htm.ProcessingResult, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}

	start := time.Now()
	defer func() {
		if ps.metricsCollector != nil {
			ps.metricsCollector.RecordProcessingTime(time.Since(start).Milliseconds())
		}
	}()

	// Stage 1: Validate input
	if err := ps.ValidateInput(ctx, input); err != nil {
		if ps.metricsCollector != nil {
			ps.metricsCollector.IncrementErrorCount()
		}
		return nil, fmt.Errorf("input validation failed: %w", err)
	}

	// Stage 2: Process the matrix
	processedData, err := ps.matrixProcessor.ProcessMatrix(ctx, input.Data)
	if err != nil {
		if ps.metricsCollector != nil {
			ps.metricsCollector.IncrementErrorCount()
		}
		return nil, fmt.Errorf("matrix processing failed: %w", err)
	}

	// Stage 3: Create result
	result := &htm.ProcessingResult{
		ID:     input.ID + "-result",
		Result: processedData,
		Metadata: htm.ResultMetadata{
			ProcessingTimeMs: time.Since(start).Milliseconds(),
			InstanceID:       ps.instanceID,
			AlgorithmVersion: ps.GetAlgorithmVersion(),
			QualityMetrics:   make(map[string]float64),
		},
		Status: htm.StatusSuccess,
	}

	// Record success metrics
	if ps.metricsCollector != nil {
		ps.metricsCollector.IncrementRequestCount()
	}

	return result, nil
}

// ValidateInput validates HTM input data without processing.
func (ps *ProcessingServiceImpl) ValidateInput(ctx context.Context, input *htm.HTMInput) error {
	if ps.validationService == nil {
		return fmt.Errorf("validation service not available")
	}

	return ps.validationService.ValidateHTMInput(input)
}

// GetSupportedVersions returns the list of supported API versions.
func (ps *ProcessingServiceImpl) GetSupportedVersions() []string {
	return []string{"v1.0"}
}

// GetAlgorithmVersion returns the current algorithm version.
func (ps *ProcessingServiceImpl) GetAlgorithmVersion() string {
	return "1.0.0"
}

// GetInstanceID returns the unique identifier for this processing instance.
func (ps *ProcessingServiceImpl) GetInstanceID() string {
	return ps.instanceID
}

// HealthCheck performs a health check on the processing service.
func (ps *ProcessingServiceImpl) HealthCheck(ctx context.Context) error {
	// Check if all required dependencies are available
	if ps.matrixProcessor == nil {
		return fmt.Errorf("matrix processor not available")
	}

	if ps.validationService == nil {
		return fmt.Errorf("validation service not available")
	}

	// Test basic processing capability with minimal data
	testInput := &htm.HTMInput{
		ID:        "health-check-test",
		Data:      [][]float64{{0.1, 0.2}, {0.3, 0.4}},
		Metadata:  htm.InputMetadata{Dimensions: []int{2, 2}, SensorID: "test", Version: "v1.0"},
		Timestamp: time.Now(),
	}

	// Quick validation test
	if err := ps.ValidateInput(ctx, testInput); err != nil {
		return fmt.Errorf("health check validation failed: %w", err)
	}

	return nil
}
