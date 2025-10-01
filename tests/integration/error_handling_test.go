package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestErrorHandlingAndValidation tests comprehensive error handling scenarios
func TestErrorHandlingAndValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	testCases := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
		retryable      bool
	}{
		{
			name:           "malformed_json",
			requestBody:    `{"invalid": json}`, // Invalid JSON
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_JSON",
			retryable:      false,
		},
		{
			name: "missing_required_fields",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"data": [][]float64{{1.0, 2.0}},
					// Missing required fields: id, metadata, timestamp
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "VALIDATION_ERROR",
			retryable:      false,
		},
		{
			name: "invalid_uuid_format",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"id":   "not-a-valid-uuid",
					"data": [][]float64{{1.0, 2.0}},
					"metadata": map[string]interface{}{
						"dimensions": []int{1, 2},
						"sensor_id":  "sensor001",
						"version":    "v1.0",
					},
					"timestamp": time.Now().Format(time.RFC3339),
				},
				"request_id": "test-invalid-uuid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_INPUT",
			retryable:      false,
		},
		{
			name: "dimension_mismatch",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"id":   "550e8400-e29b-41d4-a716-446655440000",
					"data": [][]float64{{1.0, 2.0}, {3.0, 4.0}}, // 2x2 matrix
					"metadata": map[string]interface{}{
						"dimensions": []int{3, 3}, // Claims 3x3
						"sensor_id":  "sensor001",
						"version":    "v1.0",
					},
					"timestamp": time.Now().Format(time.RFC3339),
				},
				"request_id": "test-dimension-mismatch",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_INPUT",
			retryable:      false,
		},
		{
			name: "empty_matrix_data",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"id":   "550e8400-e29b-41d4-a716-446655440001",
					"data": [][]float64{}, // Empty matrix
					"metadata": map[string]interface{}{
						"dimensions": []int{0, 0},
						"sensor_id":  "sensor001",
						"version":    "v1.0",
					},
					"timestamp": time.Now().Format(time.RFC3339),
				},
				"request_id": "test-empty-matrix",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_INPUT",
			retryable:      false,
		},
		{
			name: "inconsistent_matrix_rows",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"id":   "550e8400-e29b-41d4-a716-446655440002",
					"data": [][]float64{{1.0, 2.0, 3.0}, {4.0, 5.0}}, // Inconsistent row lengths
					"metadata": map[string]interface{}{
						"dimensions": []int{2, 3},
						"sensor_id":  "sensor001",
						"version":    "v1.0",
					},
					"timestamp": time.Now().Format(time.RFC3339),
				},
				"request_id": "test-inconsistent-matrix",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_INPUT",
			retryable:      false,
		},
		{
			name: "unsupported_api_version",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"id":   "550e8400-e29b-41d4-a716-446655440003",
					"data": [][]float64{{1.0, 2.0}},
					"metadata": map[string]interface{}{
						"dimensions": []int{1, 2},
						"sensor_id":  "sensor001",
						"version":    "v2.0", // Unsupported version
					},
					"timestamp": time.Now().Format(time.RFC3339),
				},
				"request_id": "test-unsupported-version",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "UNSUPPORTED_VERSION",
			retryable:      false,
		},
		{
			name: "invalid_sensor_id",
			requestBody: map[string]interface{}{
				"input": map[string]interface{}{
					"id":   "550e8400-e29b-41d4-a716-446655440004",
					"data": [][]float64{{1.0, 2.0}},
					"metadata": map[string]interface{}{
						"dimensions": []int{1, 2},
						"sensor_id":  "sensor@#$%", // Invalid characters
						"version":    "v1.0",
					},
					"timestamp": time.Now().Format(time.RFC3339),
				},
				"request_id": "test-invalid-sensor-id",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_INPUT",
			retryable:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if str, ok := tc.requestBody.(string); ok {
				// Handle malformed JSON case
				req, err = http.NewRequest(http.MethodPost, "/api/v1/process", strings.NewReader(str))
			} else {
				// Handle structured request body
				requestBodyBytes, marshalErr := json.Marshal(tc.requestBody)
				require.NoError(t, marshalErr)
				req, err = http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
			}

			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify status code
			assert.Equal(t, tc.expectedStatus, w.Code, "Status code mismatch for test: %s", tc.name)

			// Parse error response if possible
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)

			// For malformed JSON, response might not be parseable
			if err == nil {
				// Verify error structure
				assert.Contains(t, response, "error", "Error response should contain error field")
				assert.Contains(t, response, "request_id", "Error response should contain request_id")

				if errorObj, ok := response["error"].(map[string]interface{}); ok {
					// Verify error code
					if code, ok := errorObj["code"].(string); ok {
						assert.Equal(t, tc.expectedError, code, "Error code mismatch")
					}

					// Verify retryable flag
					if retryable, ok := errorObj["retryable"].(bool); ok {
						assert.Equal(t, tc.retryable, retryable, "Retryable flag mismatch")
					}

					// Verify error message exists
					assert.Contains(t, errorObj, "message", "Error should contain message")
					if message, ok := errorObj["message"].(string); ok {
						assert.NotEmpty(t, message, "Error message should not be empty")
					}
				}
			}
		})
	}
}

// TestErrorHandlingRecovery tests that the API can recover from errors
func TestErrorHandlingRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	// Send an invalid request
	invalidRequest := map[string]interface{}{
		"input": map[string]interface{}{
			"data": [][]float64{{}}, // Invalid empty row
		},
	}

	requestBodyBytes, err := json.Marshal(invalidRequest)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should receive error
	assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid request should return 400")

	// Now send a valid request to ensure recovery
	validRequest := map[string]interface{}{
		"input": map[string]interface{}{
			"id":   "550e8400-e29b-41d4-a716-446655440000",
			"data": [][]float64{{1.0, 2.0}},
			"metadata": map[string]interface{}{
				"dimensions": []int{1, 2},
				"sensor_id":  "sensor001",
				"version":    "v1.0",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		},
		"request_id": "recovery-test",
	}

	requestBodyBytes, err = json.Marshal(validRequest)
	require.NoError(t, err)

	req, err = http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should succeed after error recovery
	assert.Equal(t, http.StatusOK, w.Code, "Valid request should succeed after error")
}

// TestErrorHandlingConcurrent tests error handling under concurrent load
func TestErrorHandlingConcurrent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	// Test that error handling works correctly under concurrent load
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// Half valid, half invalid requests
			var requestBody map[string]interface{}

			if id%2 == 0 {
				// Valid request
				requestBody = map[string]interface{}{
					"input": map[string]interface{}{
						"id":   "550e8400-e29b-41d4-a716-446655440000",
						"data": [][]float64{{1.0, 2.0}},
						"metadata": map[string]interface{}{
							"dimensions": []int{1, 2},
							"sensor_id":  "sensor001",
							"version":    "v1.0",
						},
						"timestamp": time.Now().Format(time.RFC3339),
					},
					"request_id": fmt.Sprintf("concurrent-valid-%d", id),
				}
			} else {
				// Invalid request (missing required fields)
				requestBody = map[string]interface{}{
					"input": map[string]interface{}{
						"data": [][]float64{{1.0}},
						// Missing required fields
					},
					"request_id": fmt.Sprintf("concurrent-invalid-%d", id),
				}
			}

			requestBodyBytes, err := json.Marshal(requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify appropriate response codes
			if id%2 == 0 {
				assert.Equal(t, http.StatusOK, w.Code, "Valid concurrent request should succeed")
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid concurrent request should fail")
			}
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Helper function (reuse from htm_processing_test.go)
func setupTestRouter() *gin.Engine {
	router := gin.New()

	// TODO: These routes will fail until handlers are implemented
	// api := router.Group("/api/v1")
	// api.POST("/process", handlers.ProcessHTMInput)

	return router
}
