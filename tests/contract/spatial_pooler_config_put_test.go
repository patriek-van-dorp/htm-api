package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSpatialPoolerConfigPutEndpoint tests the PUT /api/v1/spatial-pooler/config endpoint
func TestSpatialPoolerConfigPutEndpoint(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// TODO: This will fail until the spatial pooler handler is implemented
	// router.PUT("/api/v1/spatial-pooler/config", spatialPoolerHandler.UpdateConfig)

	tests := []struct {
		name             string
		requestBody      map[string]interface{}
		expectedStatus   int
		validateResponse func(t *testing.T, response map[string]interface{})
	}{
		{
			name: "valid_config_update",
			requestBody: map[string]interface{}{
				"config": map[string]interface{}{
					"global_inhibition":  true,
					"num_active_columns": 100,
					"learning_enabled":   false,
					"boost_strength":     1.5,
					"mode":               "randomized",
				},
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "config")
				assert.Contains(t, response, "status")
				assert.Equal(t, "success", response["status"])

				config := response["config"].(map[string]interface{})

				// Verify updated values
				assert.Equal(t, true, config["global_inhibition"])
				assert.Equal(t, float64(100), config["num_active_columns"])
				assert.Equal(t, false, config["learning_enabled"])
				assert.Equal(t, 1.5, config["boost_strength"])
				assert.Equal(t, "randomized", config["mode"])
			},
		},
		{
			name: "partial_config_update",
			requestBody: map[string]interface{}{
				"config": map[string]interface{}{
					"learning_enabled": true,
					"boost_strength":   3.0,
				},
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "config")
				assert.Equal(t, "success", response["status"])

				config := response["config"].(map[string]interface{})

				// Verify updated values
				assert.Equal(t, true, config["learning_enabled"])
				assert.Equal(t, 3.0, config["boost_strength"])

				// Verify other fields are still present (not overwritten)
				assert.Contains(t, config, "input_dimensions")
				assert.Contains(t, config, "column_dimensions")
				assert.Contains(t, config, "potential_radius")
			},
		},
		{
			name: "invalid_percentage_value",
			requestBody: map[string]interface{}{
				"config": map[string]interface{}{
					"potential_pct": 1.5, // Invalid: > 1.0
				},
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "error")
				assert.Equal(t, "error", response["status"])

				error := response["error"].(map[string]interface{})
				assert.Equal(t, "configuration_error", error["type"])
				assert.Contains(t, error, "message")
				assert.Contains(t, error, "config_field")
			},
		},
		{
			name: "invalid_negative_value",
			requestBody: map[string]interface{}{
				"config": map[string]interface{}{
					"num_active_columns": -10, // Invalid: negative
				},
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "error")
				assert.Equal(t, "error", response["status"])

				error := response["error"].(map[string]interface{})
				assert.Equal(t, "configuration_error", error["type"])
			},
		},
		{
			name: "invalid_mode_value",
			requestBody: map[string]interface{}{
				"config": map[string]interface{}{
					"mode": "invalid_mode", // Invalid: not deterministic or randomized
				},
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "error")
				assert.Equal(t, "error", response["status"])

				error := response["error"].(map[string]interface{})
				assert.Equal(t, "configuration_error", error["type"])
			},
		},
		{
			name: "empty_config_object",
			requestBody: map[string]interface{}{
				"config": map[string]interface{}{},
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				// Empty config should still return current configuration
				assert.Contains(t, response, "config")
				assert.Equal(t, "success", response["status"])
			},
		},
		{
			name:        "missing_config_field",
			requestBody: map[string]interface{}{
				// Missing "config" field entirely
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "error")
				assert.Equal(t, "error", response["status"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test will currently fail since the handler is not implemented
			t.Skip("Skipping until spatial pooler handler is implemented")

			requestJSON, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("PUT", "/api/v1/spatial-pooler/config", bytes.NewBuffer(requestJSON))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			var response map[string]interface{}
			err = json.Unmarshal(recorder.Body.Bytes(), &response)
			require.NoError(t, err)

			tt.validateResponse(t, response)
		})
	}
}

// TestSpatialPoolerConfigUpdateValidation tests configuration validation rules
func TestSpatialPoolerConfigUpdateValidation(t *testing.T) {
	t.Skip("Skipping until spatial pooler handler is implemented")

	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	validationTests := []struct {
		name        string
		configField string
		value       interface{}
		shouldFail  bool
		errorType   string
	}{
		// Percentage field validations
		{"valid_potential_pct", "potential_pct", 0.5, false, ""},
		{"invalid_potential_pct_high", "potential_pct", 1.1, true, "configuration_error"},
		{"invalid_potential_pct_negative", "potential_pct", -0.1, true, "configuration_error"},

		// Integer field validations
		{"valid_num_active_columns", "num_active_columns", 50, false, ""},
		{"invalid_num_active_columns_negative", "num_active_columns", -5, true, "configuration_error"},
		{"invalid_num_active_columns_zero", "num_active_columns", 0, true, "configuration_error"},

		// Enum field validations
		{"valid_mode_deterministic", "mode", "deterministic", false, ""},
		{"valid_mode_randomized", "mode", "randomized", false, ""},
		{"invalid_mode", "mode", "invalid", true, "configuration_error"},

		// Boolean field validations
		{"valid_learning_enabled_true", "learning_enabled", true, false, ""},
		{"valid_learning_enabled_false", "learning_enabled", false, false, ""},
		{"valid_global_inhibition_true", "global_inhibition", true, false, ""},
		{"valid_global_inhibition_false", "global_inhibition", false, false, ""},

		// Numeric field validations
		{"valid_boost_strength", "boost_strength", 2.0, false, ""},
		{"invalid_boost_strength_negative", "boost_strength", -1.0, true, "configuration_error"},
		{"invalid_boost_strength_zero", "boost_strength", 0.0, true, "configuration_error"},
	}

	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := map[string]interface{}{
				"config": map[string]interface{}{
					tt.configField: tt.value,
				},
			}

			requestJSON, err := json.Marshal(requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("PUT", "/api/v1/spatial-pooler/config", bytes.NewBuffer(requestJSON))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if tt.shouldFail {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)

				var response map[string]interface{}
				err = json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Contains(t, response, "error")
				error := response["error"].(map[string]interface{})
				assert.Equal(t, tt.errorType, error["type"])
			} else {
				assert.Equal(t, http.StatusOK, recorder.Code)
			}
		})
	}
}
