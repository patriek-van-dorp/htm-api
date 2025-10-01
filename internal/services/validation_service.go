package services

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/htm-project/neural-api/internal/domain/htm"
	"github.com/htm-project/neural-api/internal/ports"
)

// ValidationServiceImpl implements the ValidationService interface.
type ValidationServiceImpl struct {
	validator *validator.Validate
	metrics   ports.MetricsCollector
}

// NewValidationService creates a new validation service.
func NewValidationService(metrics ports.MetricsCollector) ports.ValidationService {
	v := validator.New()

	// Register custom validations
	v.RegisterValidation("matrix_dimensions", validateMatrixDimensions)

	return &ValidationServiceImpl{
		validator: v,
		metrics:   metrics,
	}
}

// ValidateHTMInput validates complete HTM input structure.
func (vs *ValidationServiceImpl) ValidateHTMInput(input *htm.HTMInput) error {
	if input == nil {
		return fmt.Errorf("input cannot be nil")
	}

	// Use the input's own validation method
	if err := input.Validate(vs.validator); err != nil {
		if vs.metrics != nil {
			vs.metrics.IncrementErrorCount()
		}
		return fmt.Errorf("input validation failed: %w", err)
	}

	// Additional business rule validations
	if err := vs.validateBusinessRules(input); err != nil {
		if vs.metrics != nil {
			vs.metrics.IncrementErrorCount()
		}
		return fmt.Errorf("business rule validation failed: %w", err)
	}

	return nil
}

// ValidateMetadata validates input metadata.
func (vs *ValidationServiceImpl) ValidateMetadata(metadata *htm.InputMetadata) error {
	if metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}

	// Validate using struct tags
	if err := vs.validator.Struct(metadata); err != nil {
		return fmt.Errorf("metadata validation failed: %w", err)
	}

	// Validate shape
	if !metadata.IsValidShape() {
		return fmt.Errorf("invalid metadata shape")
	}

	return nil
}

// ValidateDimensions validates that declared dimensions match actual data.
func (vs *ValidationServiceImpl) ValidateDimensions(data [][]float64, dimensions []int) error {
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}

	if len(dimensions) < 2 {
		return fmt.Errorf("dimensions must have at least 2 values")
	}

	expectedRows, expectedCols := dimensions[0], dimensions[1]
	actualRows := len(data)

	if actualRows != expectedRows {
		return fmt.Errorf("dimension mismatch: expected %d rows, got %d", expectedRows, actualRows)
	}

	if actualRows > 0 {
		actualCols := len(data[0])
		if actualCols != expectedCols {
			return fmt.Errorf("dimension mismatch: expected %d columns, got %d", expectedCols, actualCols)
		}

		// Check all rows have same length
		for i, row := range data {
			if len(row) != expectedCols {
				return fmt.Errorf("inconsistent row length at row %d: expected %d, got %d", i, expectedCols, len(row))
			}
		}
	}

	return nil
}

// ValidateUUID validates UUID format.
func (vs *ValidationServiceImpl) ValidateUUID(uuid string) error {
	// Simple UUID validation - in real implementation would use proper UUID library
	if len(uuid) == 0 {
		return fmt.Errorf("UUID cannot be empty")
	}

	if len(uuid) != 36 {
		return fmt.Errorf("UUID must be 36 characters long")
	}

	// Check basic format (simplified)
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return fmt.Errorf("invalid UUID format")
	}

	return nil
}

// ValidateSensorID validates sensor ID format.
func (vs *ValidationServiceImpl) ValidateSensorID(sensorID string) error {
	if len(sensorID) == 0 {
		return fmt.Errorf("sensor ID cannot be empty")
	}

	if len(sensorID) > 50 {
		return fmt.Errorf("sensor ID too long: maximum 50 characters")
	}

	// Check for alphanumeric characters only
	for _, char := range sensorID {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9')) {
			return fmt.Errorf("sensor ID must contain only alphanumeric characters")
		}
	}

	return nil
}
func (vs *ValidationServiceImpl) ValidateAPIRequest(request *htm.APIRequest) error {
	if request == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate using struct tags
	if err := vs.validator.Struct(request); err != nil {
		if vs.metrics != nil {
			vs.metrics.IncrementErrorCount()
		}
		return fmt.Errorf("request validation failed: %w", err)
	}

	// Validate nested input
	if err := vs.ValidateHTMInput(&request.Input); err != nil {
		return fmt.Errorf("input validation within request failed: %w", err)
	}

	return nil
}

// validateBusinessRules validates business-specific rules for HTM input.
func (vs *ValidationServiceImpl) validateBusinessRules(input *htm.HTMInput) error {
	// Check matrix size limits for performance
	rows := len(input.Data)
	if rows > 1000 {
		return fmt.Errorf("matrix too large: %d rows exceeds maximum of 1000", rows)
	}

	if rows > 0 {
		cols := len(input.Data[0])
		if cols > 1000 {
			return fmt.Errorf("matrix too wide: %d columns exceeds maximum of 1000", cols)
		}

		// Check total elements
		totalElements := rows * cols
		if totalElements > 100000 {
			return fmt.Errorf("matrix too large: %d elements exceeds maximum of 100000", totalElements)
		}
	}

	// Validate data ranges for HTM processing - allow reasonable numeric ranges
	for i, row := range input.Data {
		for j, value := range row {
			// Check for NaN or infinite values
			if value != value { // NaN check
				return fmt.Errorf("invalid data value (NaN) at position [%d,%d]", i, j)
			}
			if value < -1000000 || value > 1000000 {
				return fmt.Errorf("data value out of reasonable range [-1000000,1000000]: %f at position [%d,%d]", value, i, j)
			}
		}
	}

	return nil
}

// Custom validation functions for validator package

// validateMatrixDimensions validates matrix dimension constraints.
func validateMatrixDimensions(fl validator.FieldLevel) bool {
	dimensions, ok := fl.Field().Interface().([]int)
	if !ok || len(dimensions) != 2 {
		return false
	}

	rows, cols := dimensions[0], dimensions[1]
	return rows > 0 && cols > 0 && rows <= 1000 && cols <= 1000
}
