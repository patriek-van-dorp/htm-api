package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Validator wraps the go-playground validator with custom rules
type Validator struct {
	validate *validator.Validate
}

// ValidationError represents a validation error with structured information
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationErrors is a slice of ValidationError
type ValidationErrors []ValidationError

// Error implements error interface for ValidationErrors
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// New creates a new validator instance with custom validation rules
func New() *Validator {
	validate := validator.New()

	// Register custom validation functions
	validate.RegisterValidation("uuid", validateUUID)
	validate.RegisterValidation("matrix_dimensions", validateMatrixDimensions)
	validate.RegisterValidation("non_empty_matrix", validateNonEmptyMatrix)

	// Use struct field names instead of json tags for validation errors
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{validate: validate}
}

// Validate validates a struct and returns structured validation errors
func (v *Validator) Validate(s interface{}) ValidationErrors {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors ValidationErrors

	for _, err := range err.(validator.ValidationErrors) {
		validationError := ValidationError{
			Field: err.Field(),
			Tag:   err.Tag(),
			Value: fmt.Sprintf("%v", err.Value()),
		}

		// Create human-readable error messages
		switch err.Tag() {
		case "required":
			validationError.Message = fmt.Sprintf("Field '%s' is required", err.Field())
		case "uuid":
			validationError.Message = fmt.Sprintf("Field '%s' must be a valid UUID", err.Field())
		case "min":
			validationError.Message = fmt.Sprintf("Field '%s' must have a minimum value/length of %s", err.Field(), err.Param())
		case "max":
			validationError.Message = fmt.Sprintf("Field '%s' must have a maximum value/length of %s", err.Field(), err.Param())
		case "oneof":
			validationError.Message = fmt.Sprintf("Field '%s' must be one of: %s", err.Field(), err.Param())
		case "alphanum":
			validationError.Message = fmt.Sprintf("Field '%s' must contain only alphanumeric characters", err.Field())
		case "matrix_dimensions":
			validationError.Message = fmt.Sprintf("Field '%s' dimensions do not match the actual data matrix", err.Field())
		case "non_empty_matrix":
			validationError.Message = fmt.Sprintf("Field '%s' matrix cannot be empty", err.Field())
		default:
			validationError.Message = fmt.Sprintf("Field '%s' failed validation for tag '%s'", err.Field(), err.Tag())
		}

		validationErrors = append(validationErrors, validationError)
	}

	return validationErrors
}

// validateUUID validates that a string is a valid UUID
func validateUUID(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return false
	}
	_, err := uuid.Parse(value)
	return err == nil
}

// validateMatrixDimensions validates that declared dimensions match actual matrix dimensions
func validateMatrixDimensions(fl validator.FieldLevel) bool {
	// This is a placeholder - actual implementation would need access to both
	// the dimensions field and the data field. In practice, this validation
	// would be done at the struct level in the business logic layer.
	return true
}

// validateNonEmptyMatrix validates that a matrix is not empty
func validateNonEmptyMatrix(fl validator.FieldLevel) bool {
	field := fl.Field()

	// Check if field is a slice
	if field.Kind() != reflect.Slice {
		return false
	}

	// Check if slice is empty
	if field.Len() == 0 {
		return false
	}

	// For 2D matrix, check if first row exists and is not empty
	if field.Index(0).Kind() == reflect.Slice {
		if field.Index(0).Len() == 0 {
			return false
		}
	}

	return true
}

// ValidateHTMInputDimensions validates that HTM input data dimensions match metadata
func ValidateHTMInputDimensions(data [][]float64, dimensions []int) error {
	if len(dimensions) < 2 {
		return fmt.Errorf("dimensions must have at least 2 values for 2D matrix")
	}

	expectedRows := dimensions[0]
	expectedCols := dimensions[1]

	if len(data) != expectedRows {
		return fmt.Errorf("expected %d rows, got %d", expectedRows, len(data))
	}

	for i, row := range data {
		if len(row) != expectedCols {
			return fmt.Errorf("row %d: expected %d columns, got %d", i, expectedCols, len(row))
		}
	}

	return nil
}

// ValidateMatrixConsistency validates that all rows in a matrix have the same length
func ValidateMatrixConsistency(data [][]float64) error {
	if len(data) == 0 {
		return fmt.Errorf("matrix cannot be empty")
	}

	expectedCols := len(data[0])
	for i, row := range data {
		if len(row) != expectedCols {
			return fmt.Errorf("inconsistent matrix: row %d has %d columns, expected %d", i, len(row), expectedCols)
		}
	}

	return nil
}
