package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProcessEndpointContract tests the contract for POST /api/v1/process
func TestProcessEndpointContract(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectResult   bool
		expectError    bool
	}{
		{
			name: "valid_htm_input_processing",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"id":   "550e8400-e29b-41d4-a716-446655440000",
					"data": [][]float64{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}},
					"metadata": map[string]interface{}{
						"dimensions": []int{2, 3},
						"sensor_id":  "sensor001",
						"version":    "v1.0",
					},
					"timestamp": "2025-09-30T10:00:00Z",
				},
				"request_id": "req-550e8400-e29b-41d4-a716-446655440001",
				"priority":   "normal",
			},
			expectedStatus: http.StatusOK,
			expectResult:   true,
			expectError:    false,
		},
		{
			name: "invalid_input_missing_required_fields",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"data": [][]float64{{1.0, 2.0}},
					// Missing id, metadata, timestamp
				},
				"request_id": "req-550e8400-e29b-41d4-a716-446655440002",
			},
			expectedStatus: http.StatusBadRequest,
			expectResult:   false,
			expectError:    true,
		},
		{
			name: "invalid_dimensions_mismatch",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"id":   "550e8400-e29b-41d4-a716-446655440003",
					"data": [][]float64{{1.0, 2.0}, {3.0, 4.0}}, // 2x2 matrix
					"metadata": map[string]interface{}{
						"dimensions": []int{2, 3}, // Claims 2x3
						"sensor_id":  "sensor001",
						"version":    "v1.0",
					},
					"timestamp": "2025-09-30T10:00:00Z",
				},
				"request_id": "req-550e8400-e29b-41d4-a716-446655440003",
			},
			expectedStatus: http.StatusBadRequest,
			expectResult:   false,
			expectError:    true,
		},
		{
			name: "empty_matrix_data",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"id":   "550e8400-e29b-41d4-a716-446655440004",
					"data": [][]float64{}, // Empty matrix
					"metadata": map[string]interface{}{
						"dimensions": []int{0, 0},
						"sensor_id":  "sensor001",
						"version":    "v1.0",
					},
					"timestamp": "2025-09-30T10:00:00Z",
				},
				"request_id": "req-550e8400-e29b-41d4-a716-446655440004",
			},
			expectedStatus: http.StatusBadRequest,
			expectResult:   false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test router with proper handlers
			router := setupTestRouter()

			// Serialize request body
			requestBodyBytes, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			// Create request
			req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request - this will fail until the route is implemented
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code, "Status code mismatch for test: %s", tt.name)

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectResult {
				require.NoError(t, err, "Failed to parse success response")

				// Verify response structure for successful processing
				assert.Contains(t, response, "result", "Success response should contain result")
				assert.Contains(t, response, "request_id", "Success response should contain request_id")
				assert.Contains(t, response, "response_time", "Success response should contain response_time")

				// Verify result structure
				if result, ok := response["result"].(map[string]interface{}); ok {
					assert.Contains(t, result, "id", "Result should contain id")
					assert.Contains(t, result, "result", "Result should contain result data")
					assert.Contains(t, result, "metadata", "Result should contain metadata")
					assert.Contains(t, result, "status", "Result should contain status")

					// Verify processing status
					assert.Equal(t, "SUCCESS", result["status"], "Status should be SUCCESS")
				}
			}

			if tt.expectError {
				// For error cases, we might not be able to parse JSON if the handler doesn't exist
				// But we should at least verify the status code is correct
				if err == nil {
					assert.Contains(t, response, "error", "Error response should contain error")
					assert.Contains(t, response, "request_id", "Error response should contain request_id")

					// Verify error structure
					if errorObj, ok := response["error"].(map[string]interface{}); ok {
						assert.Contains(t, errorObj, "code", "Error should contain code")
						assert.Contains(t, errorObj, "message", "Error should contain message")
						assert.Contains(t, errorObj, "retryable", "Error should contain retryable flag")
					}
				}
			}
		})
	}
}

// TestProcessEndpointPerformance tests that the endpoint responds within 100ms
func TestProcessEndpointPerformance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test router with proper handlers
	router := setupTestRouter()

	requestBody := map[string]interface{}{
		"input": map[string]interface{}{
			"id":   "550e8400-e29b-41d4-a716-446655440000",
			"data": [][]float64{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}},
			"metadata": map[string]interface{}{
				"dimensions": []int{2, 3},
				"sensor_id":  "sensor001",
				"version":    "v1.0",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		},
		"request_id": "perf-test-request-001",
		"priority":   "normal",
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	start := time.Now()
	router.ServeHTTP(w, req)
	duration := time.Since(start)

	// Assert response time is under 100ms (acknowledgment requirement)
	assert.Less(t, duration, 100*time.Millisecond, "Process endpoint should acknowledge within 100ms")
}
