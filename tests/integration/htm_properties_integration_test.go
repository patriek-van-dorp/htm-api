package integration

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

// HTMPropertiesTestRequest represents an HTM properties test request
type HTMPropertiesTestRequest struct {
	Input []float64 `json:"input"`
}

// TestHTMPropertiesWithActualImplementation tests HTM properties with the actual spatial pooler
func TestHTMPropertiesWithActualImplementation(t *testing.T) {
	// Setup: Create integration configuration with HTM-compliant settings
	integrationConfig := config.NewDefaultIntegrationConfig()
	integrationConfig.Application.SpatialPooler.SparsityRatio = 0.02 // 2% sparsity (HTM biological range)

	// Setup: Create spatial pooling service with actual implementation
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "htm-properties-test")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handlers
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)
	htmHandler := handlers.NewHTMValidationHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)
	api.SetupHTMValidationRoutes(router, htmHandler)

	// Phase 1: Process data to initialize spatial pooler state
	t.Run("InitializeSpatialPooler", func(t *testing.T) {
		inputSize := integrationConfig.Application.SpatialPooler.InputWidth
		testInput := make([]float64, inputSize)

		// Create a biologically realistic sparse input (5% active)
		activeInputs := int(float64(inputSize) * 0.05)
		for i := 0; i < activeInputs; i++ {
			testInput[i*20] = 1.0 // Distributed pattern
		}

		htmRequest := HTMPropertiesTestRequest{
			Input: testInput,
		}

		requestBody, err := json.Marshal(htmRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected initialization processing to succeed")

		var response map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err, "Expected valid processing response")

		// Verify: Initial processing shows HTM properties
		htmProperties := response["htm_properties"].(map[string]interface{})
		assert.True(t, htmProperties["biologically_valid"].(bool), "Expected biologically valid processing")

		sparsityAchieved := htmProperties["sparsity_achieved"].(float64)
		assert.GreaterOrEqual(t, sparsityAchieved, 0.015, "Expected sparsity >= 1.5%")
		assert.LessOrEqual(t, sparsityAchieved, 0.05, "Expected sparsity <= 5%")
	})

	// Phase 2: Validate HTM biological compliance
	t.Run("HTMBiologicalCompliance", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/validation/htm-properties", nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTM validation to succeed")

		var htmValidation map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &htmValidation)
		require.NoError(t, err, "Expected valid HTM validation response")

		// Verify: Overall HTM validity
		assert.True(t, htmValidation["valid"].(bool), "Expected HTM properties to be valid")

		// Verify: Biological compliance
		compliance := htmValidation["biological_compliance"].(map[string]interface{})
		overallScore := compliance["overall_score"].(float64)
		assert.GreaterOrEqual(t, overallScore, 0.8, "Expected high biological compliance score")

		complianceLevel := compliance["compliance_level"].(string)
		assert.Contains(t, []string{"good", "high", "excellent"}, complianceLevel, "Expected good compliance level")

		// Verify: Critical biological principles
		principles := compliance["biological_principles"].(map[string]interface{})
		assert.True(t, principles["sparse_coding"].(bool), "Expected sparse coding to be implemented")
		assert.True(t, principles["competitive_learning"].(bool), "Expected competitive learning to be implemented")
		assert.True(t, principles["hebbian_learning"].(bool), "Expected Hebbian learning to be implemented")
		assert.True(t, principles["local_inhibition"].(bool), "Expected local inhibition to be implemented")
	})

	// Phase 3: Validate sparsity properties
	t.Run("SparsityValidation", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/validation/htm-properties", nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		var htmValidation map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &htmValidation)
		require.NoError(t, err)

		// Verify: Sparsity analysis
		sparsity := htmValidation["sparsity_analysis"].(map[string]interface{})
		assert.True(t, sparsity["sparsity_valid"].(bool), "Expected sparsity validation to pass")

		currentSparsity := sparsity["current_sparsity"].(float64)
		targetSparsity := sparsity["target_sparsity"].(float64)

		// Verify: Sparsity within biological range
		assert.GreaterOrEqual(t, currentSparsity, 0.015, "Expected current sparsity >= 1.5% (biological minimum)")
		assert.LessOrEqual(t, currentSparsity, 0.05, "Expected current sparsity <= 5% (biological maximum)")

		// Verify: Target sparsity matches configuration
		assert.Equal(t, integrationConfig.Application.SpatialPooler.SparsityRatio, targetSparsity, "Expected target sparsity to match configuration")

		// Verify: Biological range is properly defined
		biologicalRange := sparsity["biological_range"].(map[string]interface{})
		minSparsity := biologicalRange["min_sparsity"].(float64)
		maxSparsity := biologicalRange["max_sparsity"].(float64)
		optimalSparsity := biologicalRange["optimal_sparsity"].(float64)

		assert.Equal(t, 0.015, minSparsity, "Expected biological minimum sparsity of 1.5%")
		assert.Equal(t, 0.05, maxSparsity, "Expected biological maximum sparsity of 5%")
		assert.GreaterOrEqual(t, optimalSparsity, minSparsity, "Expected optimal sparsity >= minimum")
		assert.LessOrEqual(t, optimalSparsity, maxSparsity, "Expected optimal sparsity <= maximum")

		// Verify: Current sparsity is within biological range
		assert.GreaterOrEqual(t, currentSparsity, minSparsity, "Expected current sparsity within biological minimum")
		assert.LessOrEqual(t, currentSparsity, maxSparsity, "Expected current sparsity within biological maximum")
	})

	// Phase 4: Validate learning and adaptation properties
	t.Run("LearningAdaptationValidation", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/validation/htm-properties", nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		var htmValidation map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &htmValidation)
		require.NoError(t, err)

		// Verify: Learning analysis
		learning := htmValidation["learning_analysis"].(map[string]interface{})
		assert.True(t, learning["learning_active"].(bool), "Expected learning to be active")

		learningRate := learning["learning_rate"].(float64)
		assert.Greater(t, learningRate, 0.0, "Expected positive learning rate")

		learningHealth := learning["learning_health"].(string)
		assert.Contains(t, []string{"good", "healthy", "excellent"}, learningHealth, "Expected good learning health")

		// Verify: Permanence updates
		permanence := learning["permanence_updates"].(map[string]interface{})
		incrementRate := permanence["increment_rate"].(float64)
		decrementRate := permanence["decrement_rate"].(float64)
		connectedThreshold := permanence["connected_threshold"].(float64)

		assert.Greater(t, incrementRate, 0.0, "Expected positive permanence increment rate")
		assert.Greater(t, decrementRate, 0.0, "Expected positive permanence decrement rate")
		assert.Greater(t, connectedThreshold, 0.0, "Expected positive connected threshold")
		assert.LessOrEqual(t, connectedThreshold, 1.0, "Expected connected threshold <= 1.0")

		// Verify: Synaptic changes
		synaptic := learning["synaptic_changes"].(map[string]interface{})
		plasticityRate := synaptic["plasticity_rate"].(float64)
		stabilityScore := synaptic["stability_score"].(float64)

		assert.GreaterOrEqual(t, plasticityRate, 0.0, "Expected non-negative plasticity rate")
		assert.GreaterOrEqual(t, stabilityScore, 0.0, "Expected non-negative stability score")
		assert.LessOrEqual(t, stabilityScore, 1.0, "Expected stability score <= 1.0")

		// Verify: Adaptation metrics
		adaptation := learning["adaptation_metrics"].(map[string]interface{})
		dutyCycleBalance := adaptation["duty_cycle_balance"].(float64)
		boostingEffectiveness := adaptation["boosting_effectiveness"].(float64)
		resourceUtilization := adaptation["resource_utilization"].(float64)

		assert.GreaterOrEqual(t, dutyCycleBalance, 0.0, "Expected non-negative duty cycle balance")
		assert.GreaterOrEqual(t, boostingEffectiveness, 0.0, "Expected non-negative boosting effectiveness")
		assert.GreaterOrEqual(t, resourceUtilization, 0.0, "Expected non-negative resource utilization")
		assert.LessOrEqual(t, resourceUtilization, 1.0, "Expected resource utilization <= 1.0")
	})

	// Phase 5: Validate inhibition and competition properties
	t.Run("InhibitionCompetitionValidation", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/validation/htm-properties", nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		var htmValidation map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &htmValidation)
		require.NoError(t, err)

		// Verify: Inhibition analysis
		inhibition := htmValidation["inhibition_analysis"].(map[string]interface{})
		assert.True(t, inhibition["inhibition_active"].(bool), "Expected inhibition to be active")

		inhibitionType := inhibition["inhibition_type"].(string)
		assert.Contains(t, []string{"global", "local", "hybrid"}, inhibitionType, "Expected valid inhibition type")

		inhibitionStrength := inhibition["inhibition_strength"].(float64)
		assert.Greater(t, inhibitionStrength, 0.0, "Expected positive inhibition strength")

		inhibitionHealth := inhibition["inhibition_health"].(string)
		assert.Contains(t, []string{"good", "healthy", "excellent"}, inhibitionHealth, "Expected good inhibition health")

		// Verify: Winner selection metrics
		winner := inhibition["winner_selection"].(map[string]interface{})
		selectionAccuracy := winner["selection_accuracy"].(float64)
		consistencyScore := winner["consistency_score"].(float64)
		competitionBalance := winner["competition_balance"].(float64)

		assert.GreaterOrEqual(t, selectionAccuracy, 0.7, "Expected good selection accuracy")
		assert.GreaterOrEqual(t, consistencyScore, 0.7, "Expected good consistency score")
		assert.GreaterOrEqual(t, competitionBalance, 0.7, "Expected good competition balance")

		winnerDistribution := winner["winner_distribution"].([]interface{})
		assert.NotEmpty(t, winnerDistribution, "Expected winner distribution data")
	})
}

// TestHTMPropertiesSparsityThresholds tests sparsity threshold handling
func TestHTMPropertiesSparsityThresholds(t *testing.T) {
	// Test Case 1: Very sparse input (<0.5% active)
	t.Run("VerySparseinput", func(t *testing.T) {
		integrationConfig := config.NewDefaultIntegrationConfig()
		integrationConfig.Application.SpatialPooler.SparsityRatio = 0.02

		spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "sparse-test")
		require.NoError(t, err)

		spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		api.SetupSpatialPoolerRoutes(router, spatialHandler)

		// Create very sparse input (0.3% active)
		inputSize := integrationConfig.Application.SpatialPooler.InputWidth
		testInput := make([]float64, inputSize)
		activeInputs := int(float64(inputSize) * 0.003) // 0.3% active
		for i := 0; i < activeInputs; i++ {
			testInput[i*100] = 1.0
		}

		htmRequest := HTMPropertiesTestRequest{
			Input: testInput,
		}

		requestBody, err := json.Marshal(htmRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Should handle sparse input gracefully
		assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, recorder.Code,
			"Expected either successful processing or graceful rejection of very sparse input")
	})

	// Test Case 2: Dense input (>10% active)
	t.Run("DenseInput", func(t *testing.T) {
		integrationConfig := config.NewDefaultIntegrationConfig()
		integrationConfig.Application.SpatialPooler.SparsityRatio = 0.02

		spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "dense-test")
		require.NoError(t, err)

		spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		api.SetupSpatialPoolerRoutes(router, spatialHandler)

		// Create dense input (15% active)
		inputSize := integrationConfig.Application.SpatialPooler.InputWidth
		testInput := make([]float64, inputSize)
		activeInputs := int(float64(inputSize) * 0.15) // 15% active
		for i := 0; i < activeInputs; i++ {
			testInput[i] = 1.0
		}

		htmRequest := HTMPropertiesTestRequest{
			Input: testInput,
		}

		requestBody, err := json.Marshal(htmRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Should handle dense input gracefully
		assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, recorder.Code,
			"Expected either successful processing or graceful rejection of very dense input")

		if recorder.Code == http.StatusOK {
			var response map[string]interface{}
			err = json.Unmarshal(recorder.Body.Bytes(), &response)
			require.NoError(t, err)

			// If processed successfully, output should still maintain sparsity
			result := response["result"].(map[string]interface{})
			sparsity := result["sparsity"].(float64)
			assert.LessOrEqual(t, sparsity, 0.05, "Expected output sparsity <= 5% even with dense input")
		}
	})
}

// TestHTMPropertiesConfigurationValidation tests HTM configuration validation
func TestHTMPropertiesConfigurationValidation(t *testing.T) {
	// Test: Invalid sparsity configuration
	t.Run("InvalidSparsityConfiguration", func(t *testing.T) {
		integrationConfig := config.NewDefaultIntegrationConfig()

		// Create invalid configuration with sparsity outside biological range
		invalidConfig := &htm.SpatialPoolerConfig{
			InputWidth:          1024,
			ColumnCount:         2048,
			SparsityRatio:       0.15, // Invalid: too high (>10%)
			Mode:                htm.SpatialPoolerModeDeterministic,
			LearningEnabled:     true,
			LearningRate:        0.1,
			MaxBoost:            3.0,
			BoostStrength:       0.5,
			InhibitionRadius:    16,
			LocalAreaDensity:    0.02,
			MinOverlapThreshold: 5,
			MaxProcessingTimeMs: 10,
			SemanticThresholds: htm.SemanticThresholds{
				SimilarInputMinOverlap:   0.5,
				DifferentInputMaxOverlap: 0.1,
			},
		}

		spatialService, err := services.NewSpatialPoolingService(invalidConfig, "invalid-config-test")

		// Should either reject invalid configuration or handle it gracefully
		if err != nil {
			// Expected: Configuration validation should catch invalid sparsity
			assert.Contains(t, err.Error(), "sparsity", "Expected error message to mention sparsity")
		} else {
			// If service is created, it should handle the invalid configuration
			configHandler := handlers.NewSpatialPoolerConfigHandler(spatialService)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			api.SetupSpatialPoolerRoutes(router, configHandler)

			// The system should detect and report the invalid configuration
			req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/status", nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code, "Expected status endpoint to respond")
		}
	})

	// Test: Valid HTM configuration
	t.Run("ValidHTMConfiguration", func(t *testing.T) {
		integrationConfig := config.NewDefaultIntegrationConfig()

		// Ensure configuration is within HTM biological ranges
		assert.GreaterOrEqual(t, integrationConfig.Application.SpatialPooler.SparsityRatio, 0.015,
			"Expected sparsity ratio >= 1.5%")
		assert.LessOrEqual(t, integrationConfig.Application.SpatialPooler.SparsityRatio, 0.05,
			"Expected sparsity ratio <= 5%")

		spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "valid-config-test")
		require.NoError(t, err, "Expected valid HTM configuration to be accepted")

		htmHandler := handlers.NewHTMValidationHandler(spatialService)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		api.SetupHTMValidationRoutes(router, htmHandler)

		req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/validation/htm-properties", nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTM validation to succeed with valid configuration")

		var htmValidation map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &htmValidation)
		require.NoError(t, err)

		assert.True(t, htmValidation["valid"].(bool), "Expected HTM properties to be valid with proper configuration")
	})
}
