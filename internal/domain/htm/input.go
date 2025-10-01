package htm

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// HTMInput represents multi-dimensional arrays containing spatial-temporal pattern data
// for HTM neural network processing.
type HTMInput struct {
	ID        string        `json:"id" validate:"required,uuid"`
	Data      [][]float64   `json:"data" validate:"required,min=1"`
	Metadata  InputMetadata `json:"metadata" validate:"required"`
	Timestamp time.Time     `json:"timestamp" validate:"required"`
}

// Validate validates the HTMInput struct and its business rules
func (h *HTMInput) Validate(validator *validator.Validate) error {
	// First run standard validation
	if err := validator.Struct(h); err != nil {
		return err
	}

	// Custom business logic validation
	if err := h.validateMatrixConsistency(); err != nil {
		return err
	}

	if err := h.validateDimensionsMatch(); err != nil {
		return err
	}

	return nil
}

// validateMatrixConsistency ensures all rows have the same length
func (h *HTMInput) validateMatrixConsistency() error {
	if len(h.Data) == 0 {
		return &ValidationError{Field: "data", Message: "matrix cannot be empty"}
	}

	expectedCols := len(h.Data[0])
	for i, row := range h.Data {
		if len(row) != expectedCols {
			return &ValidationError{
				Field:   "data",
				Message: fmt.Sprintf("inconsistent matrix: row %d has %d columns, expected %d", i, len(row), expectedCols),
			}
		}
	}

	return nil
}

// validateDimensionsMatch ensures metadata dimensions match actual data dimensions
func (h *HTMInput) validateDimensionsMatch() error {
	if len(h.Metadata.Dimensions) < 2 {
		return &ValidationError{
			Field:   "metadata.dimensions",
			Message: "dimensions must have at least 2 values for 2D matrix",
		}
	}

	expectedRows := h.Metadata.Dimensions[0]
	expectedCols := h.Metadata.Dimensions[1]

	if len(h.Data) != expectedRows {
		return &ValidationError{
			Field:   "metadata.dimensions",
			Message: fmt.Sprintf("expected %d rows, got %d", expectedRows, len(h.Data)),
		}
	}

	if len(h.Data) > 0 && len(h.Data[0]) != expectedCols {
		return &ValidationError{
			Field:   "metadata.dimensions",
			Message: fmt.Sprintf("expected %d columns, got %d", expectedCols, len(h.Data[0])),
		}
	}

	return nil
}

// GetDimensions returns the actual dimensions of the data matrix
func (h *HTMInput) GetDimensions() []int {
	if len(h.Data) == 0 {
		return []int{0, 0}
	}
	return []int{len(h.Data), len(h.Data[0])}
}

// IsEmpty returns true if the HTMInput contains no data
func (h *HTMInput) IsEmpty() bool {
	return len(h.Data) == 0 || (len(h.Data) == 1 && len(h.Data[0]) == 0)
}

// ValidationError represents a domain-specific validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
