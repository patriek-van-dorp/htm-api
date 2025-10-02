package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/api"
	"github.com/htm-project/neural-api/internal/domain/htm"
	"github.com/htm-project/neural-api/internal/handlers"
	"github.com/htm-project/neural-api/internal/infrastructure/config"
	"github.com/htm-project/neural-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SpatialPoolerConfigUpdateRequest represents the configuration update request
type SpatialPoolerConfigUpdateRequest struct {
	Configuration *htm.SpatialPoolerConfig `json:"configuration" validate:"required"`
	ApplyMode     string                   `json:"apply_mode,omitempty"`
	ValidateOnly  bool                     `json:"validate_only,omitempty"`
}

// SpatialPoolerConfigUpdateResponse represents the configuration update response
type SpatialPoolerConfigUpdateResponse struct {
	Success          bool                     `json:"success"`
	Timestamp        string                   `json:"timestamp"`
	Applied          bool                     `json:"applied"`
	ValidationResult ConfigValidationResult   `json:"validation_result"`
	PreviousConfig   *htm.SpatialPoolerConfig `json:"previous_config,omitempty"`
	NewConfig        *htm.SpatialPoolerConfig `json:"new_config"`
	Impact           ConfigurationImpact      `json:"impact"`
}

// ConfigValidationResult represents validation results
type ConfigValidationResult struct {
	Valid     bool                `json:"valid"`
	Errors    []string            `json:"errors,omitempty"`
	Warnings  []string            `json:"warnings,omitempty"`
	HTMValid  bool                `json:"htm_valid"`
	HTMChecks HTMValidationChecks `json:"htm_checks"`
}

// HTMValidationChecks represents HTM-specific validation checks
type HTMValidationChecks struct {
	SparsityRatioValid      bool `json:"sparsity_ratio_valid"`
	ColumnCountValid        bool `json:"column_count_valid"`
	InputWidthValid         bool `json:"input_width_valid"`
	PermanenceValid         bool `json:"permanence_valid"`
	InhibitionValid         bool `json:"inhibition_valid"`
	LearningParametersValid bool `json:"learning_parameters_valid"`
}

// ConfigurationImpact represents the impact of configuration changes
type ConfigurationImpact struct {
	RequiresReinitialization bool     `json:"requires_reinitialization"`
	RequiresLearningReset    bool     `json:"requires_learning_reset"`
	AffectedComponents       []string `json:"affected_components"`
	PerformanceImpact        string   `json:"performance_impact"`
	MemoryImpact             string   `json:"memory_impact"`
}

// TestSpatialPoolerConfigUpdateIntegrationValid tests valid configuration update
func TestSpatialPoolerConfigUpdateIntegrationValid(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create spatial pooler configuration handler
	configHandler := handlers.NewSpatialPoolerConfigHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, configHandler)

	// Prepare: Create valid configuration update
	newConfig := &htm.SpatialPoolerConfig{
		InputWidth:              2048,  // Changed from default
		ColumnCount:             4096,  // Changed from default
		PotentialRadius:         128,   // Changed from default
		PotentialPct:            0.9,   // Changed from default
		ConnectedPermanence:     0.35,  // Changed from default
		PermanenceIncrement:     0.06,  // Changed from default
		PermanenceDecrement:     0.01,  // Changed from default
		MinOverlapThreshold:     8,     // Changed from default
		LocalAreaDensity:        0.025, // Changed from default
		SparsityRatio:           0.025, // Changed from default
		StimulusThreshold:       2,     // Changed from default
		GlobalInhibition:        true,
		DutyCyclePeriod:         120,   // Changed from default
		BoostStrength:           2.5,   // Changed from default
		SynPermInactiveDec:      0.01,  // Changed from default
		SynPermActiveInc:        0.06,  // Changed from default
		SynPermConnected:        0.35,  // Changed from default
		MinPctOverlapDutyCycles: 0.002, // Changed from default
		MinPctActiveDutyCycles:  0.002, // Changed from default
		WrapAround:              false,
	}

	updateRequest := SpatialPoolerConfigUpdateRequest{
		Configuration: newConfig,
		ApplyMode:     "immediate",
		ValidateOnly:  false,
	}

	// Execute: Make configuration update request
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", "/api/v1/spatial-pooler/config", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected successful configuration update")

	// Verify: Response body structure
	var response SpatialPoolerConfigUpdateResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Update success
	assert.True(t, response.Success, "Expected configuration update to succeed")
	assert.True(t, response.Applied, "Expected configuration to be applied")
	assert.NotEmpty(t, response.Timestamp, "Expected timestamp to be present")

	// Verify: Validation result
	validation := response.ValidationResult
	assert.True(t, validation.Valid, "Expected configuration to be valid")
	assert.Empty(t, validation.Errors, "Expected no validation errors")
	assert.True(t, validation.HTMValid, "Expected HTM validation to pass")

	// Verify: HTM validation checks
	htmChecks := validation.HTMChecks
	assert.True(t, htmChecks.SparsityRatioValid, "Expected sparsity ratio to be valid")
	assert.True(t, htmChecks.ColumnCountValid, "Expected column count to be valid")
	assert.True(t, htmChecks.InputWidthValid, "Expected input width to be valid")
	assert.True(t, htmChecks.PermanenceValid, "Expected permanence values to be valid")
	assert.True(t, htmChecks.InhibitionValid, "Expected inhibition parameters to be valid")
	assert.True(t, htmChecks.LearningParametersValid, "Expected learning parameters to be valid")

	// Verify: Previous and new configuration
	assert.NotNil(t, response.PreviousConfig, "Expected previous configuration to be returned")
	assert.NotNil(t, response.NewConfig, "Expected new configuration to be returned")
	assert.Equal(t, newConfig.InputWidth, response.NewConfig.InputWidth, "Expected input width to be updated")
	assert.Equal(t, newConfig.ColumnCount, response.NewConfig.ColumnCount, "Expected column count to be updated")
	assert.Equal(t, newConfig.SparsityRatio, response.NewConfig.SparsityRatio, "Expected sparsity ratio to be updated")

	// Verify: Configuration impact
	impact := response.Impact
	assert.True(t, impact.RequiresReinitialization, "Expected reinitialization to be required for major changes")
	assert.NotEmpty(t, impact.AffectedComponents, "Expected affected components to be listed")
	assert.Contains(t, []string{"minimal", "moderate", "significant"}, impact.PerformanceImpact, "Expected valid performance impact")
	assert.Contains(t, []string{"minimal", "moderate", "significant"}, impact.MemoryImpact, "Expected valid memory impact")
}

// TestSpatialPoolerConfigUpdateIntegrationValidationOnly tests validation-only mode
func TestSpatialPoolerConfigUpdateIntegrationValidationOnly(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	configHandler := handlers.NewSpatialPoolerConfigHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, configHandler)

	// Prepare: Create configuration for validation only
	newConfig := &htm.SpatialPoolerConfig{
		InputWidth:              1024,
		ColumnCount:             2048,
		PotentialRadius:         64,
		PotentialPct:            0.8,
		ConnectedPermanence:     0.3,
		PermanenceIncrement:     0.05,
		PermanenceDecrement:     0.008,
		MinOverlapThreshold:     5,
		LocalAreaDensity:        0.02,
		SparsityRatio:           0.03, // Slightly different from default
		StimulusThreshold:       1,
		GlobalInhibition:        true,
		DutyCyclePeriod:         100,
		BoostStrength:           2.0,
		SynPermInactiveDec:      0.008,
		SynPermActiveInc:        0.05,
		SynPermConnected:        0.3,
		MinPctOverlapDutyCycles: 0.001,
		MinPctActiveDutyCycles:  0.001,
		WrapAround:              false,
	}

	updateRequest := SpatialPoolerConfigUpdateRequest{
		Configuration: newConfig,
		ApplyMode:     "validation",
		ValidateOnly:  true,
	}

	// Execute: Make validation request
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", "/api/v1/spatial-pooler/config", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected successful validation")

	// Verify: Response body
	var response SpatialPoolerConfigUpdateResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Validation success but not applied
	assert.True(t, response.Success, "Expected validation to succeed")
	assert.False(t, response.Applied, "Expected configuration NOT to be applied in validation-only mode")

	// Verify: Validation results
	assert.True(t, response.ValidationResult.Valid, "Expected configuration to be valid")
	assert.True(t, response.ValidationResult.HTMValid, "Expected HTM validation to pass")

	// Verify: Impact analysis still provided
	assert.NotNil(t, response.Impact, "Expected impact analysis even in validation-only mode")
}

// TestSpatialPoolerConfigUpdateIntegrationInvalidSparsity tests invalid sparsity configuration
func TestSpatialPoolerConfigUpdateIntegrationInvalidSparsity(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	configHandler := handlers.NewSpatialPoolerConfigHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, configHandler)

	// Prepare: Create configuration with invalid sparsity
	invalidConfig := &htm.SpatialPoolerConfig{
		InputWidth:              1024,
		ColumnCount:             2048,
		PotentialRadius:         64,
		PotentialPct:            0.8,
		ConnectedPermanence:     0.3,
		PermanenceIncrement:     0.05,
		PermanenceDecrement:     0.008,
		MinOverlapThreshold:     5,
		LocalAreaDensity:        0.02,
		SparsityRatio:           0.15, // Invalid: too high (>10%)
		StimulusThreshold:       1,
		GlobalInhibition:        true,
		DutyCyclePeriod:         100,
		BoostStrength:           2.0,
		SynPermInactiveDec:      0.008,
		SynPermActiveInc:        0.05,
		SynPermConnected:        0.3,
		MinPctOverlapDutyCycles: 0.001,
		MinPctActiveDutyCycles:  0.001,
		WrapAround:              false,
	}

	updateRequest := SpatialPoolerConfigUpdateRequest{
		Configuration: invalidConfig,
		ApplyMode:     "immediate",
		ValidateOnly:  false,
	}

	// Execute: Make configuration update request
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", "/api/v1/spatial-pooler/config", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status (should be 400 for invalid configuration)
	assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected bad request for invalid sparsity")

	// Verify: Response body
	var response SpatialPoolerConfigUpdateResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Update failure
	assert.False(t, response.Success, "Expected configuration update to fail")
	assert.False(t, response.Applied, "Expected configuration NOT to be applied")

	// Verify: Validation errors
	validation := response.ValidationResult
	assert.False(t, validation.Valid, "Expected configuration to be invalid")
	assert.NotEmpty(t, validation.Errors, "Expected validation errors")
	assert.False(t, validation.HTMValid, "Expected HTM validation to fail")

	// Verify: Specific HTM validation failures
	htmChecks := validation.HTMChecks
	assert.False(t, htmChecks.SparsityRatioValid, "Expected sparsity ratio validation to fail")

	// Verify: Error messages mention sparsity
	errorMessages := validation.Errors
	sparsityErrorFound := false
	for _, msg := range errorMessages {
		if contains(msg, "sparsity") || contains(msg, "sparse") {
			sparsityErrorFound = true
			break
		}
	}
	assert.True(t, sparsityErrorFound, "Expected sparsity error message")
}

// TestSpatialPoolerConfigUpdateIntegrationInvalidPermanence tests invalid permanence values
func TestSpatialPoolerConfigUpdateIntegrationInvalidPermanence(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	configHandler := handlers.NewSpatialPoolerConfigHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, configHandler)

	// Prepare: Create configuration with invalid permanence values
	invalidConfig := &htm.SpatialPoolerConfig{
		InputWidth:              1024,
		ColumnCount:             2048,
		PotentialRadius:         64,
		PotentialPct:            0.8,
		ConnectedPermanence:     1.5,    // Invalid: > 1.0
		PermanenceIncrement:     -0.05,  // Invalid: negative
		PermanenceDecrement:     -0.008, // Invalid: negative
		MinOverlapThreshold:     5,
		LocalAreaDensity:        0.02,
		SparsityRatio:           0.02,
		StimulusThreshold:       1,
		GlobalInhibition:        true,
		DutyCyclePeriod:         100,
		BoostStrength:           2.0,
		SynPermInactiveDec:      -0.008, // Invalid: negative
		SynPermActiveInc:        -0.05,  // Invalid: negative
		SynPermConnected:        1.5,    // Invalid: > 1.0
		MinPctOverlapDutyCycles: 0.001,
		MinPctActiveDutyCycles:  0.001,
		WrapAround:              false,
	}

	updateRequest := SpatialPoolerConfigUpdateRequest{
		Configuration: invalidConfig,
		ApplyMode:     "immediate",
		ValidateOnly:  false,
	}

	// Execute: Make configuration update request
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", "/api/v1/spatial-pooler/config", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected bad request for invalid permanence")

	// Verify: Response body
	var response SpatialPoolerConfigUpdateResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Validation failure
	assert.False(t, response.Success, "Expected configuration update to fail")
	assert.False(t, response.ValidationResult.Valid, "Expected configuration to be invalid")
	assert.False(t, response.ValidationResult.HTMValid, "Expected HTM validation to fail")
	assert.False(t, response.ValidationResult.HTMChecks.PermanenceValid, "Expected permanence validation to fail")

	// Verify: Permanence error messages
	errorMessages := response.ValidationResult.Errors
	permanenceErrorFound := false
	for _, msg := range errorMessages {
		if contains(msg, "permanence") {
			permanenceErrorFound = true
			break
		}
	}
	assert.True(t, permanenceErrorFound, "Expected permanence error message")
}

// TestSpatialPoolerConfigUpdateIntegrationHTMProperties tests HTM property validation
func TestSpatialPoolerConfigUpdateIntegrationHTMProperties(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	configHandler := handlers.NewSpatialPoolerConfigHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, configHandler)

	// Prepare: Create configuration that maintains HTM properties
	htmValidConfig := &htm.SpatialPoolerConfig{
		InputWidth:              1024,
		ColumnCount:             2048,
		PotentialRadius:         64,
		PotentialPct:            0.8,
		ConnectedPermanence:     0.3,
		PermanenceIncrement:     0.05,
		PermanenceDecrement:     0.008,
		MinOverlapThreshold:     5,
		LocalAreaDensity:        0.02,
		SparsityRatio:           0.025, // Within HTM biological range
		StimulusThreshold:       1,
		GlobalInhibition:        true,
		DutyCyclePeriod:         100,
		BoostStrength:           2.0,
		SynPermInactiveDec:      0.008,
		SynPermActiveInc:        0.05,
		SynPermConnected:        0.3,
		MinPctOverlapDutyCycles: 0.001,
		MinPctActiveDutyCycles:  0.001,
		WrapAround:              false,
	}

	updateRequest := SpatialPoolerConfigUpdateRequest{
		Configuration: htmValidConfig,
		ApplyMode:     "immediate",
		ValidateOnly:  false,
	}

	// Execute: Make configuration update request
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", "/api/v1/spatial-pooler/config", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected successful HTM configuration update")

	// Verify: HTM validation passes
	var response SpatialPoolerConfigUpdateResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	assert.True(t, response.Success, "Expected HTM configuration update to succeed")
	assert.True(t, response.ValidationResult.HTMValid, "Expected HTM validation to pass")

	// Verify: All HTM checks pass
	htmChecks := response.ValidationResult.HTMChecks
	assert.True(t, htmChecks.SparsityRatioValid, "Expected sparsity ratio to be HTM valid")
	assert.True(t, htmChecks.ColumnCountValid, "Expected column count to be HTM valid")
	assert.True(t, htmChecks.InputWidthValid, "Expected input width to be HTM valid")
	assert.True(t, htmChecks.PermanenceValid, "Expected permanence values to be HTM valid")
	assert.True(t, htmChecks.InhibitionValid, "Expected inhibition to be HTM valid")
	assert.True(t, htmChecks.LearningParametersValid, "Expected learning parameters to be HTM valid")
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsInternal(s, substr))))
}

func containsInternal(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
