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

// TestResponseTimeRequirements tests that the API meets performance requirements
func TestResponseTimeRequirements(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	testCases := []struct {
		name               string
		matrixSize         string
		data               [][]float64
		maxAcknowledgeTime time.Duration
		maxProcessingTime  time.Duration
	}{
		{
			name:               "small_matrix_2x2",
			matrixSize:         "2x2",
			data:               [][]float64{{1.0, 2.0}, {3.0, 4.0}},
			maxAcknowledgeTime: 100 * time.Millisecond,
			maxProcessingTime:  500 * time.Millisecond,
		},
		{
			name:               "medium_matrix_5x5",
			matrixSize:         "5x5",
			data:               generateTestMatrix(5, 5),
			maxAcknowledgeTime: 100 * time.Millisecond,
			maxProcessingTime:  1 * time.Second,
		},
		{
			name:               "large_matrix_10x10",
			matrixSize:         "10x10",
			data:               generateTestMatrix(10, 10),
			maxAcknowledgeTime: 100 * time.Millisecond,
			maxProcessingTime:  2 * time.Second,
		},
		{
			name:               "very_large_matrix_20x20",
			matrixSize:         "20x20",
			data:               generateTestMatrix(20, 20),
			maxAcknowledgeTime: 100 * time.Millisecond,
			maxProcessingTime:  5 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requestBody := map[string]interface{}{
				"input": map[string]interface{}{
					"id":   fmt.Sprintf("perf-test-%s", tc.matrixSize),
					"data": tc.data,
					"metadata": map[string]interface{}{
						"dimensions": []int{len(tc.data), len(tc.data[0])},
						"sensor_id":  "performance_sensor",
						"version":    "v1.0",
					},
					"timestamp": time.Now().Format(time.RFC3339),
				},
				"request_id": fmt.Sprintf("perf-request-%s", tc.matrixSize),
				"priority":   "normal",
			}

			// Measure total response time
			start := time.Now()

			requestBodyBytes, err := json.Marshal(requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			totalResponseTime := time.Since(start)

			// Critical requirement: Acknowledgment within 100ms
			assert.Less(t, totalResponseTime, tc.maxAcknowledgeTime,
				"API should acknowledge %s request within %v", tc.matrixSize, tc.maxAcknowledgeTime)

			// Verify response status
			assert.Equal(t, http.StatusOK, w.Code, "Performance test should succeed for %s", tc.matrixSize)

			// Parse response to verify processing time metadata
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)

			if err == nil {
				if result, ok := response["result"].(map[string]interface{}); ok {
					if metadata, ok := result["metadata"].(map[string]interface{}); ok {
						if processingTimeMs, ok := metadata["processing_time_ms"].(float64); ok {
							processingTime := time.Duration(processingTimeMs) * time.Millisecond
							assert.Less(t, processingTime, tc.maxProcessingTime,
								"Processing time for %s should be under %v", tc.matrixSize, tc.maxProcessingTime)
						}
					}
				}
			}
		})
	}
}

// TestThroughputRequirements tests sustained throughput under load
func TestThroughputRequirements(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	// Test parameters
	duration := 10 * time.Second
	targetThroughput := 50.0 // requests per second

	var requestCount int64
	var successCount int64
	var mu sync.Mutex

	endTime := time.Now().Add(duration)
	var wg sync.WaitGroup

	// Launch multiple workers to generate load
	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for time.Now().Before(endTime) {
				requestBody := map[string]interface{}{
					"input": map[string]interface{}{
						"id":   fmt.Sprintf("throughput-test-%d-%d", workerID, time.Now().UnixNano()),
						"data": [][]float64{{1.0, 2.0}, {3.0, 4.0}},
						"metadata": map[string]interface{}{
							"dimensions": []int{2, 2},
							"sensor_id":  fmt.Sprintf("throughput_sensor_%d", workerID),
							"version":    "v1.0",
						},
						"timestamp": time.Now().Format(time.RFC3339),
					},
					"request_id": fmt.Sprintf("throughput-request-%d-%d", workerID, time.Now().UnixNano()),
				}

				requestBodyBytes, err := json.Marshal(requestBody)
				if err != nil {
					continue
				}

				req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
				if err != nil {
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				mu.Lock()
				requestCount++
				if w.Code == http.StatusOK {
					successCount++
				}
				mu.Unlock()

				// Small delay to prevent overwhelming the system
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()

	// Calculate actual throughput
	actualThroughput := float64(requestCount) / duration.Seconds()
	successRate := float64(successCount) / float64(requestCount)

	// Assertions
	assert.Greater(t, actualThroughput, targetThroughput,
		"Should achieve at least %.1f requests/second, got %.1f", targetThroughput, actualThroughput)
	assert.Greater(t, successRate, 0.95,
		"Should have at least 95%% success rate, got %.2f", successRate)

	t.Logf("Throughput test results: %.1f req/s, %.2f%% success rate", actualThroughput, successRate*100)
}

// TestLatencyDistribution tests response time distribution
func TestLatencyDistribution(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	numRequests := 100
	responseTimes := make([]time.Duration, 0, numRequests)

	for i := 0; i < numRequests; i++ {
		requestBody := map[string]interface{}{
			"input": map[string]interface{}{
				"id":   fmt.Sprintf("latency-test-%d", i),
				"data": [][]float64{{float64(i), float64(i + 1)}, {float64(i + 2), float64(i + 3)}},
				"metadata": map[string]interface{}{
					"dimensions": []int{2, 2},
					"sensor_id":  "latency_sensor",
					"version":    "v1.0",
				},
				"timestamp": time.Now().Format(time.RFC3339),
			},
			"request_id": fmt.Sprintf("latency-request-%d", i),
		}

		start := time.Now()

		requestBodyBytes, err := json.Marshal(requestBody)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseTime := time.Since(start)
		responseTimes = append(responseTimes, responseTime)

		// Each request should succeed
		assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i)
	}

	// Calculate percentiles
	p50, p95, p99 := calculatePercentiles(responseTimes)
	avg := calculateAverageResponseTime(responseTimes)

	// Performance requirements
	assert.Less(t, p50, 50*time.Millisecond, "P50 latency should be under 50ms")
	assert.Less(t, p95, 100*time.Millisecond, "P95 latency should be under 100ms")
	assert.Less(t, p99, 200*time.Millisecond, "P99 latency should be under 200ms")
	assert.Less(t, avg, 75*time.Millisecond, "Average latency should be under 75ms")

	t.Logf("Latency distribution - Avg: %v, P50: %v, P95: %v, P99: %v", avg, p50, p95, p99)
}

// TestMemoryUsageUnderLoad tests memory efficiency during sustained load
func TestMemoryUsageUnderLoad(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the full application is implemented
	router := setupTestRouter()

	// Generate load for a period and monitor that response times don't degrade
	duration := 5 * time.Second
	endTime := time.Now().Add(duration)

	var responseTimes []time.Duration
	var mu sync.Mutex

	var wg sync.WaitGroup
	numWorkers := 5

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			requestCounter := 0
			for time.Now().Before(endTime) {
				requestBody := map[string]interface{}{
					"input": map[string]interface{}{
						"id":   fmt.Sprintf("memory-test-%d-%d", workerID, requestCounter),
						"data": generateTestMatrix(5, 5), // Medium size matrix
						"metadata": map[string]interface{}{
							"dimensions": []int{5, 5},
							"sensor_id":  fmt.Sprintf("memory_sensor_%d", workerID),
							"version":    "v1.0",
						},
						"timestamp": time.Now().Format(time.RFC3339),
					},
					"request_id": fmt.Sprintf("memory-request-%d-%d", workerID, requestCounter),
				}

				start := time.Now()

				requestBodyBytes, err := json.Marshal(requestBody)
				if err != nil {
					continue
				}

				req, err := http.NewRequest(http.MethodPost, "/api/v1/process", bytes.NewBuffer(requestBodyBytes))
				if err != nil {
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				responseTime := time.Since(start)

				mu.Lock()
				responseTimes = append(responseTimes, responseTime)
				mu.Unlock()

				requestCounter++
				time.Sleep(50 * time.Millisecond) // Prevent overwhelming
			}
		}(i)
	}

	wg.Wait()

	// Verify that response times remain consistent (no memory leaks causing degradation)
	if len(responseTimes) > 20 {
		firstQuarter := responseTimes[:len(responseTimes)/4]
		lastQuarter := responseTimes[3*len(responseTimes)/4:]

		avgFirst := calculateAverageResponseTime(firstQuarter)
		avgLast := calculateAverageResponseTime(lastQuarter)

		// Response time shouldn't degrade by more than 50% over time
		degradationRatio := float64(avgLast) / float64(avgFirst)
		assert.Less(t, degradationRatio, 1.5,
			"Response time shouldn't degrade significantly (got %.2fx degradation)", degradationRatio)

		t.Logf("Memory usage test - First quarter avg: %v, Last quarter avg: %v", avgFirst, avgLast)
	}
}

// Helper functions

func calculatePercentiles(times []time.Duration) (p50, p95, p99 time.Duration) {
	if len(times) == 0 {
		return 0, 0, 0
	}

	// Simple sorting and percentile calculation
	sorted := make([]time.Duration, len(times))
	copy(sorted, times)

	// Bubble sort for simplicity (not efficient but adequate for tests)
	for i := 0; i < len(sorted); i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	p50Index := int(float64(len(sorted)) * 0.5)
	p95Index := int(float64(len(sorted)) * 0.95)
	p99Index := int(float64(len(sorted)) * 0.99)

	if p50Index >= len(sorted) {
		p50Index = len(sorted) - 1
	}
	if p95Index >= len(sorted) {
		p95Index = len(sorted) - 1
	}
	if p99Index >= len(sorted) {
		p99Index = len(sorted) - 1
	}

	return sorted[p50Index], sorted[p95Index], sorted[p99Index]
}
