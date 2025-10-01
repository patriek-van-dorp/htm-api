package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConcurrentRequestHandling tests handling of multiple simultaneous requests
func TestConcurrentRequestHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	// Test configuration
	numConcurrentRequests := 20
	requestsPerGoroutine := 5

	var wg sync.WaitGroup
	results := make(chan TestResult, numConcurrentRequests*requestsPerGoroutine)

	// Launch concurrent goroutines
	for i := 0; i < numConcurrentRequests; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < requestsPerGoroutine; j++ {
				requestID := fmt.Sprintf("concurrent-%d-%d", goroutineID, j)

				requestBody := map[string]interface{}{
					"input": map[string]interface{}{
						"id":   "550e8400-e29b-41d4-a716-446655440000",
						"data": generateVariableMatrix(goroutineID, j),
						"metadata": map[string]interface{}{
							"dimensions": []int{2, 3},
							"sensor_id":  fmt.Sprintf("sensor%03d", goroutineID),
							"version":    "v1.0",
						},
						"timestamp": time.Now().Format(time.RFC3339),
					},
					"request_id": requestID,
					"priority":   "normal",
				}

				result := makeTimedRequest(router, requestBody, requestID)
				results <- result
			}
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	close(results)

	// Analyze results
	var successCount, errorCount int
	var responseTimes []time.Duration

	for result := range results {
		if result.Success {
			successCount++
		} else {
			errorCount++
		}
		responseTimes = append(responseTimes, result.ResponseTime)
	}

	totalRequests := numConcurrentRequests * requestsPerGoroutine

	// Assertions
	assert.Equal(t, totalRequests, successCount+errorCount, "All requests should be accounted for")
	assert.GreaterOrEqual(t, successCount, totalRequests/2, "At least half of requests should succeed")

	// Performance assertions (these will fail until implementation is optimized)
	avgResponseTime := calculateAverageResponseTime(responseTimes)
	assert.Less(t, avgResponseTime, 200*time.Millisecond, "Average response time should be under 200ms")

	// No response should take longer than 1 second
	for _, rt := range responseTimes {
		assert.Less(t, rt, 1*time.Second, "No request should take longer than 1 second")
	}
}

// TestConcurrentRequestIsolation tests that concurrent requests don't interfere with each other
func TestConcurrentRequestIsolation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	// Test that concurrent requests with different data don't interfere
	numGoroutines := 10
	var wg sync.WaitGroup
	results := make(chan IsolationTestResult, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Each goroutine uses unique input data
			uniqueData := generateUniqueMatrix(id)

			requestBody := map[string]interface{}{
				"input": map[string]interface{}{
					"id":   fmt.Sprintf("isolation-test-%d", id),
					"data": uniqueData,
					"metadata": map[string]interface{}{
						"dimensions": []int{len(uniqueData), len(uniqueData[0])},
						"sensor_id":  fmt.Sprintf("isolation_sensor_%d", id),
						"version":    "v1.0",
					},
					"timestamp": time.Now().Format(time.RFC3339),
				},
				"request_id": fmt.Sprintf("isolation-request-%d", id),
			}

			requestBodyBytes, err := json.Marshal(requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			result := IsolationTestResult{
				GoroutineID: id,
				StatusCode:  w.Code,
				InputData:   uniqueData,
			}

			// Parse response if successful
			if w.Code == http.StatusOK {
				var response map[string]interface{}
				if json.Unmarshal(w.Body.Bytes(), &response) == nil {
					if resultObj, ok := response["result"].(map[string]interface{}); ok {
						if resultData, ok := resultObj["result"].([]interface{}); ok {
							result.OutputData = resultData
						}
					}
				}
			}

			results <- result
		}(i)
	}

	wg.Wait()
	close(results)

	// Verify isolation - each request should have received appropriate response for its input
	for result := range results {
		assert.Equal(t, http.StatusOK, result.StatusCode, "Request %d should succeed", result.GoroutineID)

		if result.OutputData != nil {
			// Verify output dimensions match input dimensions
			assert.Len(t, result.OutputData, len(result.InputData), "Output should have same number of rows as input")
		}
	}
}

// TestConcurrentRequestResourceManagement tests resource usage under load
func TestConcurrentRequestResourceManagement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	// Test high load scenario
	numGoroutines := 50
	requestsPerGoroutine := 10

	startTime := time.Now()
	var wg sync.WaitGroup
	successChannel := make(chan bool, numGoroutines*requestsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < requestsPerGoroutine; j++ {
				requestBody := map[string]interface{}{
					"input": map[string]interface{}{
						"id":   fmt.Sprintf("load-test-%d-%d", id, j),
						"data": [][]float64{{float64(id), float64(j)}, {float64(j), float64(id)}},
						"metadata": map[string]interface{}{
							"dimensions": []int{2, 2},
							"sensor_id":  fmt.Sprintf("load_sensor_%d", id),
							"version":    "v1.0",
						},
						"timestamp": time.Now().Format(time.RFC3339),
					},
					"request_id": fmt.Sprintf("load-request-%d-%d", id, j),
				}

				requestBodyBytes, err := json.Marshal(requestBody)
				if err != nil {
					successChannel <- false
					continue
				}

				req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
				if err != nil {
					successChannel <- false
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Consider any response (even errors) as handled successfully from resource perspective
				successChannel <- w.Code != 0
			}
		}(i)
	}

	wg.Wait()
	close(successChannel)

	duration := time.Since(startTime)

	// Count successful responses
	handledRequests := 0
	for handled := range successChannel {
		if handled {
			handledRequests++
		}
	}

	totalRequests := numGoroutines * requestsPerGoroutine

	// Resource management assertions
	assert.Equal(t, totalRequests, handledRequests, "All requests should be handled")
	assert.Less(t, duration, 30*time.Second, "High load test should complete within 30 seconds")

	// Calculate throughput
	throughput := float64(handledRequests) / duration.Seconds()
	assert.Greater(t, throughput, 10.0, "Should handle at least 10 requests per second")
}

// Helper types and functions

type TestResult struct {
	RequestID    string
	Success      bool
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

type IsolationTestResult struct {
	GoroutineID int
	StatusCode  int
	InputData   [][]float64
	OutputData  []interface{}
}

func makeTimedRequest(router *gin.Engine, requestBody map[string]interface{}, requestID string) TestResult {
	start := time.Now()

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return TestResult{
			RequestID:    requestID,
			Success:      false,
			ResponseTime: time.Since(start),
			Error:        err,
		}
	}

	req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return TestResult{
			RequestID:    requestID,
			Success:      false,
			ResponseTime: time.Since(start),
			Error:        err,
		}
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	responseTime := time.Since(start)

	return TestResult{
		RequestID:    requestID,
		Success:      w.Code == http.StatusOK,
		StatusCode:   w.Code,
		ResponseTime: responseTime,
	}
}

func generateVariableMatrix(goroutineID, requestID int) [][]float64 {
	// Generate matrix with values based on IDs to ensure uniqueness
	base := float64(goroutineID*100 + requestID)
	return [][]float64{
		{base + 1, base + 2, base + 3},
		{base + 4, base + 5, base + 6},
	}
}

func generateUniqueMatrix(id int) [][]float64 {
	// Generate a unique matrix for each ID
	size := 2 + (id % 3) // Vary size between 2x2, 3x3, 4x4
	matrix := make([][]float64, size)

	for i := range matrix {
		matrix[i] = make([]float64, size)
		for j := range matrix[i] {
			matrix[i][j] = float64(id*100 + i*10 + j)
		}
	}

	return matrix
}
