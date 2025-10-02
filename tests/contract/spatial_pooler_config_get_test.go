package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSpatialPoolerConfigGetEndpoint tests the GET /api/v1/spatial-pooler/config endpoint
func TestSpatialPoolerConfigGetEndpoint(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// TODO: This will fail until the spatial pooler handler is implemented
	// router.GET("/api/v1/spatial-pooler/config", spatialPoolerHandler.GetConfig)

	t.Run("get_current_config", func(t *testing.T) {
		// This test will currently fail since the handler is not implemented
		t.Skip("Skipping until spatial pooler handler is implemented")

		req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/config", nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err)

		// Validate response structure
		assert.Contains(t, response, "config")
		assert.Contains(t, response, "status")
		assert.Equal(t, "success", response["status"])

		config := response["config"].(map[string]interface{})

		// Validate required configuration fields
		requiredFields := []string{
			"input_dimensions",
			"column_dimensions",
			"potential_radius",
			"potential_pct",
			"global_inhibition",
			"num_active_columns",
			"stimulus_threshold",
			"syn_perm_inc",
			"syn_perm_dec",
			"syn_perm_connected",
			"syn_perm_trim_threshold",
			"duty_cycle_period",
			"boost_strength",
			"learning_enabled",
			"mode",
		}

		for _, field := range requiredFields {
			assert.Contains(t, config, field, "Config should contain field: %s", field)
		}

		// Validate field types and ranges
		if inputDims, ok := config["input_dimensions"].([]interface{}); ok {
			assert.NotEmpty(t, inputDims, "input_dimensions should not be empty")
			for _, dim := range inputDims {
				if dimFloat, ok := dim.(float64); ok {
					assert.Greater(t, dimFloat, 0.0, "input_dimensions should be positive")
				}
			}
		}

		if columnDims, ok := config["column_dimensions"].([]interface{}); ok {
			assert.NotEmpty(t, columnDims, "column_dimensions should not be empty")
			for _, dim := range columnDims {
				if dimFloat, ok := dim.(float64); ok {
					assert.Greater(t, dimFloat, 0.0, "column_dimensions should be positive")
				}
			}
		}

		// Validate percentage fields are between 0 and 1
		percentageFields := []string{"potential_radius", "potential_pct", "syn_perm_inc", "syn_perm_dec", "syn_perm_connected", "syn_perm_trim_threshold"}
		for _, field := range percentageFields {
			if value, ok := config[field].(float64); ok {
				assert.GreaterOrEqual(t, value, 0.0, "%s should be >= 0", field)
				assert.LessOrEqual(t, value, 1.0, "%s should be <= 1", field)
			}
		}

		// Validate integer fields are positive
		integerFields := []string{"num_active_columns", "stimulus_threshold", "duty_cycle_period"}
		for _, field := range integerFields {
			if value, ok := config[field].(float64); ok {
				assert.GreaterOrEqual(t, value, 0.0, "%s should be >= 0", field)
			}
		}

		// Validate boost_strength is positive
		if boostStrength, ok := config["boost_strength"].(float64); ok {
			assert.Greater(t, boostStrength, 0.0, "boost_strength should be positive")
		}

		// Validate mode is valid enum value
		if mode, ok := config["mode"].(string); ok {
			assert.Contains(t, []string{"deterministic", "randomized"}, mode, "mode should be 'deterministic' or 'randomized'")
		}

		// Validate learning_enabled is boolean
		if learningEnabled, ok := config["learning_enabled"]; ok {
			_, isBool := learningEnabled.(bool)
			assert.True(t, isBool, "learning_enabled should be boolean")
		}

		// Validate global_inhibition is boolean
		if globalInhibition, ok := config["global_inhibition"]; ok {
			_, isBool := globalInhibition.(bool)
			assert.True(t, isBool, "global_inhibition should be boolean")
		}
	})
}

// TestSpatialPoolerConfigDefaultValues tests that the default configuration values are sensible
func TestSpatialPoolerConfigDefaultValues(t *testing.T) {
	t.Skip("Skipping until spatial pooler handler is implemented")

	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/config", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	config := response["config"].(map[string]interface{})

	// Test reasonable default values for HTM spatial pooling
	expectedDefaults := map[string]interface{}{
		"potential_radius":        0.1,             // 10% radius is reasonable
		"potential_pct":           0.5,             // 50% potential connections
		"syn_perm_connected":      0.1,             // Standard HTM value
		"syn_perm_inc":            0.05,            // Standard increment
		"syn_perm_dec":            0.008,           // Standard decrement
		"syn_perm_trim_threshold": 0.01,            // Standard trim threshold
		"duty_cycle_period":       1000,            // Standard period
		"boost_strength":          2.0,             // Standard boost
		"learning_enabled":        true,            // Learning enabled by default
		"mode":                    "deterministic", // Deterministic by default
	}

	for field, expectedValue := range expectedDefaults {
		if actualValue, exists := config[field]; exists {
			assert.Equal(t, expectedValue, actualValue, "Default value for %s should be %v", field, expectedValue)
		}
	}

	// Validate that num_active_columns is reasonable (2-5% of total columns)
	if numActiveColumns, ok := config["num_active_columns"].(float64); ok {
		if columnDims, ok := config["column_dimensions"].([]interface{}); ok && len(columnDims) > 0 {
			totalColumns := 1.0
			for _, dim := range columnDims {
				if dimFloat, ok := dim.(float64); ok {
					totalColumns *= dimFloat
				}
			}
			sparsity := numActiveColumns / totalColumns
			assert.GreaterOrEqual(t, sparsity, 0.02, "Default sparsity should be >= 2%")
			assert.LessOrEqual(t, sparsity, 0.05, "Default sparsity should be <= 5%")
		}
	}
}
