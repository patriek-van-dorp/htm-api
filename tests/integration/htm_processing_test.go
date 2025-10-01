package integration

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

// TestHTMProcessingWorkflow tests the complete HTM input processing workflow
func TestHTMProcessingWorkflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	// Test case: Complete workflow from input to processed output
	testCases := []struct {
		name        string
		inputData   [][]float64
		dimensions  []int
		sensorID    string
		expectError bool
	}{
		{
			name:        "basic_2x3_matrix_processing",
			inputData:   [][]float64{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}},
			dimensions:  []int{2, 3},
			sensorID:    "sensor001",
			expectError: false,
		},
		{
			name:        "large_matrix_processing",
			inputData:   generateTestMatrix(10, 10),
			dimensions:  []int{10, 10},
			sensorID:    "sensor002",
			expectError: false,
		},
		{
			name:        "single_row_matrix",
			inputData:   [][]float64{{1.0, 2.0, 3.0, 4.0, 5.0}},
			dimensions:  []int{1, 5},
			sensorID:    "sensor003",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Step 1: Prepare HTM input
			requestBody := map[string]interface{}{
				"input": map[string]interface{}{
					"id":   "550e8400-e29b-41d4-a716-446655440000",
					"data": tc.inputData,
					"metadata": map[string]interface{}{
						"dimensions": tc.dimensions,
						"sensor_id":  tc.sensorID,
						"version":    "v1.0",
					},
					"timestamp": time.Now().Format(time.RFC3339),
				},
				"request_id": "workflow-test-" + tc.name,
				"priority":   "normal",
			}

			// Step 2: Send processing request
			requestBodyBytes, err := json.Marshal(requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if tc.expectError {
				assert.Equal(t, http.StatusBadRequest, w.Code, "Expected error response")
				return
			}

			// Step 3: Verify successful processing
			assert.Equal(t, http.StatusOK, w.Code, "Processing should succeed")

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			// Step 4: Verify response structure and content
			assert.Contains(t, response, "result", "Response should contain result")
			assert.Contains(t, response, "request_id", "Response should contain request_id")

			if result, ok := response["result"].(map[string]interface{}); ok {
				// Verify result data matches input dimensions
				if resultData, ok := result["result"].([]interface{}); ok {
					assert.Len(t, resultData, len(tc.inputData), "Result should have same number of rows as input")

					// Verify each row has correct number of columns
					for i, row := range resultData {
						if rowData, ok := row.([]interface{}); ok {
							assert.Len(t, rowData, len(tc.inputData[i]), "Result row %d should have same columns as input", i)
						}
					}
				}

				// Verify processing metadata
				if metadata, ok := result["metadata"].(map[string]interface{}); ok {
					assert.Contains(t, metadata, "processing_time_ms", "Metadata should contain processing time")
					assert.Contains(t, metadata, "instance_id", "Metadata should contain instance ID")
					assert.Contains(t, metadata, "algorithm_version", "Metadata should contain algorithm version")
				}

				// Verify processing status
				assert.Equal(t, "SUCCESS", result["status"], "Processing status should be SUCCESS")
			}
		})
	}
}

// TestHTMProcessingChaining tests API chaining capability
func TestHTMProcessingChaining(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	// Step 1: Process initial input
	firstInput := map[string]interface{}{
		"input": map[string]interface{}{
			"id":   "chain-input-001",
			"data": [][]float64{{1.0, 2.0}, {3.0, 4.0}},
			"metadata": map[string]interface{}{
				"dimensions": []int{2, 2},
				"sensor_id":  "sensor_chain",
				"version":    "v1.0",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		},
		"request_id": "chain-request-001",
	}

	firstResponse := makeProcessingRequest(t, router, firstInput)
	require.NotNil(t, firstResponse, "First processing should succeed")

	// Step 2: Use first response as input to second processing
	if firstResult, ok := firstResponse["result"].(map[string]interface{}); ok {
		if firstResultData, ok := firstResult["result"]; ok {
			secondInput := map[string]interface{}{
				"input": map[string]interface{}{
					"id":   "chain-input-002",
					"data": firstResultData, // Chain the output
					"metadata": map[string]interface{}{
						"dimensions": []int{2, 2},
						"sensor_id":  "sensor_chain",
						"version":    "v1.0",
					},
					"timestamp": time.Now().Format(time.RFC3339),
				},
				"request_id": "chain-request-002",
			}

			secondResponse := makeProcessingRequest(t, router, secondInput)
			require.NotNil(t, secondResponse, "Second processing should succeed")

			// Verify that chaining worked
			assert.Contains(t, secondResponse, "result", "Chained processing should succeed")
		}
	}
}

// TestHTMProcessingDataIntegrity tests that data integrity is maintained
func TestHTMProcessingDataIntegrity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	// Test with known input data and verify output maintains expected properties
	inputData := [][]float64{{1.0, 0.0, 1.0}, {0.0, 1.0, 0.0}, {1.0, 1.0, 0.0}}

	requestBody := map[string]interface{}{
		"input": map[string]interface{}{
			"id":   "integrity-test-001",
			"data": inputData,
			"metadata": map[string]interface{}{
				"dimensions": []int{3, 3},
				"sensor_id":  "integrity_sensor",
				"version":    "v1.0",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		},
		"request_id": "integrity-request-001",
	}

	response := makeProcessingRequest(t, router, requestBody)
	require.NotNil(t, response, "Processing should succeed")

	// Verify data integrity properties
	if result, ok := response["result"].(map[string]interface{}); ok {
		if resultData, ok := result["result"].([]interface{}); ok {
			// Verify dimensions are preserved
			assert.Len(t, resultData, 3, "Output should have same number of rows")

			for i, row := range resultData {
				if rowData, ok := row.([]interface{}); ok {
					assert.Len(t, rowData, 3, "Row %d should have same number of columns", i)

					// Verify all values are valid numbers
					for j, val := range rowData {
						assert.IsType(t, float64(0), val, "Value at [%d][%d] should be a number", i, j)
					}
				}
			}
		}
	}
}

// Helper functions

func makeProcessingRequest(t *testing.T, router *gin.Engine, requestBody map[string]interface{}) map[string]interface{} {
	requestBodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return nil
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response
}
