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

// TestSpatialPoolerProcessEndpoint tests the POST /api/v1/spatial-pooler/process endpoint
func TestSpatialPoolerProcessEndpoint(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// TODO: This will fail until the spatial pooler handler is implemented
	// router.POST("/api/v1/spatial-pooler/process", spatialPoolerHandler.Process)

	tests := []struct {
		name             string
		requestBody      map[string]interface{}
		expectedStatus   int
		validateResponse func(t *testing.T, response map[string]interface{})
	}{
		{
			name: "valid_spatial_pooler_request",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"encoder_output": map[string]interface{}{
						"width":       2048,
						"active_bits": []int{10, 25, 67, 89, 134, 256, 445, 678, 789, 901},
						"sparsity":    0.0048828125,
					},
					"input_id":         "test-uuid-001",
					"learning_enabled": true,
					"metadata": map[string]interface{}{
						"sensor_type":     "categorical",
						"encoder_version": "1.0.0",
					},
				},
				"config": map[string]interface{}{
					"mode":             "deterministic",
					"learning_enabled": true,
				},
				"request_id": "req-test-001",
				"client_id":  "test-client",
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				// Validate response structure
				assert.Contains(t, response, "result")
				assert.Contains(t, response, "request_id")
				assert.Contains(t, response, "status")
				assert.Equal(t, "success", response["status"])

				result := response["result"].(map[string]interface{})
				assert.Contains(t, result, "sdr")
				assert.Contains(t, result, "input_id")
				assert.Contains(t, result, "processing_time_ms")
				assert.Contains(t, result, "sparsity_level")

				// Validate SDR structure
				sdr := result["sdr"].(map[string]interface{})
				assert.Contains(t, sdr, "width")
				assert.Contains(t, sdr, "active_bits")
				assert.Contains(t, sdr, "sparsity")

				// Validate sparsity is in 2-5% range
				sparsity := result["sparsity_level"].(float64)
				assert.GreaterOrEqual(t, sparsity, 0.02, "Sparsity should be >= 2%")
				assert.LessOrEqual(t, sparsity, 0.05, "Sparsity should be <= 5%")

				// Validate processing time is under 10ms
				processingTime := result["processing_time_ms"].(float64)
				assert.LessOrEqual(t, processingTime, 10.0, "Processing time should be <= 10ms")
			},
		},
		{
			name: "invalid_input_oversized",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"encoder_output": map[string]interface{}{
						"width":       4096, // Exceeds expected 2048
						"active_bits": []int{10, 25, 67},
						"sparsity":    0.001,
					},
					"input_id": "test-uuid-002",
				},
				"request_id": "req-test-002",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "error")
				assert.Contains(t, response, "status")
				assert.Equal(t, "error", response["status"])

				error := response["error"].(map[string]interface{})
				assert.Contains(t, error, "type")
				assert.Contains(t, error, "message")
				assert.Equal(t, "invalid_input", error["type"])
			},
		},
		{
			name: "missing_required_fields",
			requestBody: map[string]interface{}{
				"request_id": "req-test-003",
				// Missing input field
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

			req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestJSON))
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

// TestSpatialPoolerProcessDeterministicBehavior tests that identical inputs produce identical outputs in deterministic mode
func TestSpatialPoolerProcessDeterministicBehavior(t *testing.T) {
	t.Skip("Skipping until spatial pooler handler is implemented")

	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	requestBody := map[string]interface{}{
		"input": map[string]interface{}{
			"encoder_output": map[string]interface{}{
				"width":       2048,
				"active_bits": []int{100, 200, 300, 400, 500},
				"sparsity":    0.00244,
			},
			"input_id":         "deterministic-test",
			"learning_enabled": false, // Disable learning for deterministic test
		},
		"config": map[string]interface{}{
			"mode":             "deterministic",
			"learning_enabled": false,
		},
		"request_id": "deterministic-req",
	}

	requestJSON, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Make first request
	req1, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestJSON))
	require.NoError(t, err)
	req1.Header.Set("Content-Type", "application/json")

	recorder1 := httptest.NewRecorder()
	router.ServeHTTP(recorder1, req1)

	// Make second request with identical input
	req2, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestJSON))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")

	recorder2 := httptest.NewRecorder()
	router.ServeHTTP(recorder2, req2)

	// Both should succeed
	assert.Equal(t, http.StatusOK, recorder1.Code)
	assert.Equal(t, http.StatusOK, recorder2.Code)

	var response1, response2 map[string]interface{}
	err = json.Unmarshal(recorder1.Body.Bytes(), &response1)
	require.NoError(t, err)
	err = json.Unmarshal(recorder2.Body.Bytes(), &response2)
	require.NoError(t, err)

	// Extract active_bits from both responses
	result1 := response1["result"].(map[string]interface{})
	result2 := response2["result"].(map[string]interface{})

	sdr1 := result1["sdr"].(map[string]interface{})
	sdr2 := result2["sdr"].(map[string]interface{})

	activeBits1 := sdr1["active_bits"]
	activeBits2 := sdr2["active_bits"]

	// In deterministic mode with learning disabled, active bits should be identical
	assert.Equal(t, activeBits1, activeBits2, "Deterministic mode should produce identical outputs for identical inputs")
}
