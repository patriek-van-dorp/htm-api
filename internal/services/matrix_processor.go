package services

import (
	"context"
	"fmt"
	"time"

	"github.com/htm-project/neural-api/internal/ports"
)

// MatrixProcessorImpl implements the MatrixProcessor interface.
type MatrixProcessorImpl struct {
	metrics ports.MetricsCollector
}

// NewMatrixProcessor creates a new matrix processor.
func NewMatrixProcessor(metrics ports.MetricsCollector) ports.MatrixProcessor {
	return &MatrixProcessorImpl{
		metrics: metrics,
	}
}

// ProcessMatrix performs matrix processing operations on 2D data.
func (mp *MatrixProcessorImpl) ProcessMatrix(ctx context.Context, data [][]float64) ([][]float64, error) {
	if data == nil {
		return nil, fmt.Errorf("data cannot be nil")
	}

	start := time.Now()
	defer func() {
		if mp.metrics != nil {
			mp.metrics.RecordProcessingTime(time.Since(start).Milliseconds())
		}
	}()

	// Validate matrix
	if err := mp.ValidateMatrix(data); err != nil {
		return nil, fmt.Errorf("matrix validation failed: %w", err)
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	rows := len(data)
	cols := len(data[0])
	result := make([][]float64, rows)

	// HTM processing simulation
	for i := 0; i < rows; i++ {
		result[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			originalValue := data[i][j]

			// Simulate HTM spatial pooling
			spatialValue := mp.applySpatialPooling(originalValue, i, j)

			// Simulate HTM temporal memory
			temporalValue := mp.applyTemporalMemory(spatialValue, i, j)

			result[i][j] = temporalValue
		}

		// Check context cancellation periodically
		if i%100 == 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
		}
	}

	// Record metrics
	if mp.metrics != nil {
		mp.metrics.IncrementRequestCount()
	}

	return result, nil
}

// ValidateMatrix validates that a matrix is suitable for processing.
func (mp *MatrixProcessorImpl) ValidateMatrix(data [][]float64) error {
	if data == nil {
		return fmt.Errorf("matrix cannot be nil")
	}

	rows := len(data)

	// Check dimensions
	if rows == 0 {
		return fmt.Errorf("matrix cannot be empty")
	}

	cols := len(data[0])
	if cols == 0 {
		return fmt.Errorf("matrix rows cannot be empty")
	}

	// Check maximum size limits (prevent memory issues)
	maxSize := 10000 // Maximum elements
	if rows*cols > maxSize {
		return fmt.Errorf("matrix too large: %dx%d exceeds maximum size of %d elements", rows, cols, maxSize)
	}

	// Check consistency and validate data ranges
	for i, row := range data {
		if len(row) != cols {
			return fmt.Errorf("inconsistent row lengths: row %d has %d columns, expected %d", i, len(row), cols)
		}

		for j, value := range row {
			// Check for NaN or Inf values
			if value != value { // NaN check
				return fmt.Errorf("matrix contains NaN at position [%d,%d]", i, j)
			}
			if value > 1e308 || value < -1e308 { // Inf check
				return fmt.Errorf("matrix contains Inf at position [%d,%d]", i, j)
			}

			// Allow reasonable numeric ranges for general processing
			if value < -1000000 || value > 1000000 {
				return fmt.Errorf("matrix values out of reasonable range [-1000000, 1000000], found %f at position [%d,%d]", value, i, j)
			}
		}
	}

	return nil
}

// GetMatrixDimensions returns the dimensions of a matrix.
func (mp *MatrixProcessorImpl) GetMatrixDimensions(data [][]float64) (rows, cols int) {
	if data == nil || len(data) == 0 {
		return 0, 0
	}
	return len(data), len(data[0])
}

// IsMatrixConsistent checks if all rows have the same length.
func (mp *MatrixProcessorImpl) IsMatrixConsistent(data [][]float64) bool {
	if data == nil || len(data) == 0 {
		return true
	}

	expectedCols := len(data[0])
	for _, row := range data {
		if len(row) != expectedCols {
			return false
		}
	}
	return true
}

// NormalizeMatrix normalizes matrix values to a standard range.
func (mp *MatrixProcessorImpl) NormalizeMatrix(data [][]float64) ([][]float64, error) {
	if data == nil {
		return nil, fmt.Errorf("data cannot be nil")
	}

	if !mp.IsMatrixConsistent(data) {
		return nil, fmt.Errorf("matrix is not consistent")
	}

	rows, cols := mp.GetMatrixDimensions(data)
	result := mp.CreateEmptyMatrix(rows, cols)

	// Find min and max values
	var min, max float64
	if rows > 0 && cols > 0 {
		min = data[0][0]
		max = data[0][0]

		for _, row := range data {
			for _, value := range row {
				if value < min {
					min = value
				}
				if value > max {
					max = value
				}
			}
		}
	}

	// Normalize to [0, 1] range
	rangeVal := max - min
	if rangeVal == 0 {
		// All values are the same
		for i := range result {
			for j := range result[i] {
				result[i][j] = 0.5
			}
		}
	} else {
		for i, row := range data {
			for j, value := range row {
				result[i][j] = (value - min) / rangeVal
			}
		}
	}

	return result, nil
}

// CreateEmptyMatrix creates a matrix with specified dimensions.
func (mp *MatrixProcessorImpl) CreateEmptyMatrix(rows, cols int) [][]float64 {
	if rows <= 0 || cols <= 0 {
		return nil
	}

	result := make([][]float64, rows)
	for i := range result {
		result[i] = make([]float64, cols)
	}
	return result
}

// applySpatialPooling simulates HTM spatial pooling algorithm.
func (mp *MatrixProcessorImpl) applySpatialPooling(value float64, row, col int) float64 {
	// Simplified spatial pooling simulation
	// In real HTM, this would involve complex sparse distributed representations

	// Apply some spatial context based on neighboring positions
	spatialFactor := 0.1 * float64((row+col)%10) / 10.0
	sparsity := 0.02 // 2% sparsity typical for HTM

	// Threshold-based activation
	if value+spatialFactor > 0.5 {
		return 1.0 * sparsity // Active column
	}
	return 0.0 // Inactive column
}

// applyTemporalMemory simulates HTM temporal memory algorithm.
func (mp *MatrixProcessorImpl) applyTemporalMemory(spatialValue float64, row, col int) float64 {
	// Simplified temporal memory simulation
	// In real HTM, this would maintain sequence memory and predictions

	if spatialValue > 0 {
		// Apply temporal context
		temporalFactor := 0.8 // Learning rate
		prediction := 0.2     // Base prediction strength

		return spatialValue*temporalFactor + prediction
	}

	return spatialValue
}
