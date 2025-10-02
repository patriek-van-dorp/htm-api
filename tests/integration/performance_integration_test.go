package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/api"
	"github.com/htm-project/neural-api/internal/handlers"
	"github.com/htm-project/neural-api/internal/infrastructure/config"
	"github.com/htm-project/neural-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PerformanceTestRequest represents a performance test request
type PerformanceTestRequest struct {
	Input []float64 `json:"input"`
}

// TestPerformanceRequirements tests the <100ms response time requirement
func TestPerformanceRequirements(t *testing.T) {
	// Setup: Create integration configuration optimized for performance
	integrationConfig := config.NewDefaultIntegrationConfig()
	integrationConfig.Performance.ResponseTimeTargetMs = 100

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "perf-test")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Test: Multiple requests to ensure consistent performance
	const numRequests = 50
	var totalTime time.Duration
	var maxTime time.Duration
	successCount := 0

	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)
	// Create a realistic sparse pattern
	for i := 0; i < inputSize/30; i++ {
		testInput[i*30] = 1.0
	}

	performanceRequest := PerformanceTestRequest{
		Input: testInput,
	}

	requestBody, err := json.Marshal(performanceRequest)
	require.NoError(t, err)

	for i := 0; i < numRequests; i++ {
		// Measure request processing time
		startTime := time.Now()

		req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		endTime := time.Now()
		requestTime := endTime.Sub(startTime)

		if recorder.Code == http.StatusOK {
			successCount++
			totalTime += requestTime
			if requestTime > maxTime {
				maxTime = requestTime
			}

			// Verify: Individual request under 100ms
			assert.LessOrEqual(t, requestTime.Milliseconds(), int64(100), "Expected individual request under 100ms")
		}
	}

	// Verify: All requests succeeded
	assert.Equal(t, numRequests, successCount, "Expected all performance test requests to succeed")

	// Verify: Average response time under target
	avgTime := totalTime / time.Duration(successCount)
	assert.LessOrEqual(t, avgTime.Milliseconds(), int64(100), "Expected average response time under 100ms")

	// Verify: Maximum response time reasonable
	assert.LessOrEqual(t, maxTime.Milliseconds(), int64(150), "Expected maximum response time under 150ms")

	t.Logf("Performance Results: Avg: %v, Max: %v, Success Rate: %d/%d",
		avgTime, maxTime, successCount, numRequests)
}

// TestPerformanceUnderLoad tests performance with sustained load
func TestPerformanceUnderLoad(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "load-test")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Test: Sustained load for 10 seconds
	const testDuration = 10 * time.Second
	const targetRPS = 20 // requests per second

	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)
	for i := 0; i < inputSize/20; i++ {
		testInput[i*20] = 1.0
	}

	performanceRequest := PerformanceTestRequest{
		Input: testInput,
	}

	requestBody, err := json.Marshal(performanceRequest)
	require.NoError(t, err)

	startTime := time.Now()
	requestCount := 0
	successCount := 0
	var totalProcessingTime time.Duration

	for time.Since(startTime) < testDuration {
		requestStart := time.Now()

		req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		requestEnd := time.Now()
		processingTime := requestEnd.Sub(requestStart)

		requestCount++
		if recorder.Code == http.StatusOK {
			successCount++
			totalProcessingTime += processingTime

			// Verify: Each request still meets performance target under load
			assert.LessOrEqual(t, processingTime.Milliseconds(), int64(120), "Expected request under 120ms even under load")
		}

		// Control request rate
		time.Sleep(time.Second / time.Duration(targetRPS))
	}

	actualDuration := time.Since(startTime)
	actualRPS := float64(requestCount) / actualDuration.Seconds()
	successRate := float64(successCount) / float64(requestCount)
	avgProcessingTime := totalProcessingTime / time.Duration(successCount)

	// Verify: Sustained throughput
	assert.GreaterOrEqual(t, actualRPS, float64(targetRPS)*0.9, "Expected to maintain 90% of target RPS under load")

	// Verify: High success rate under load
	assert.GreaterOrEqual(t, successRate, 0.95, "Expected 95% success rate under sustained load")

	// Verify: Average processing time still reasonable
	assert.LessOrEqual(t, avgProcessingTime.Milliseconds(), int64(100), "Expected average processing time under 100ms under load")

	t.Logf("Load Test Results: RPS: %.2f, Success Rate: %.2f%%, Avg Time: %v",
		actualRPS, successRate*100, avgProcessingTime)
}

// TestPerformanceLargeInput tests performance with large input datasets
func TestPerformanceLargeInput(t *testing.T) {
	// Setup: Create integration configuration with larger input
	integrationConfig := config.NewDefaultIntegrationConfig()
	integrationConfig.Application.SpatialPooler.InputWidth = 4096  // Large input
	integrationConfig.Application.SpatialPooler.ColumnCount = 8192 // Large column count

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "large-test")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Test: Large input processing performance
	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)

	// Create a realistic sparse pattern for large input
	activeInputs := int(float64(inputSize) * 0.02) // 2% active
	for i := 0; i < activeInputs; i++ {
		testInput[i*int(inputSize/activeInputs)] = 1.0
	}

	performanceRequest := PerformanceTestRequest{
		Input: testInput,
	}

	requestBody, err := json.Marshal(performanceRequest)
	require.NoError(t, err)

	// Execute: Multiple requests with large input
	const numLargeRequests = 10
	var totalTime time.Duration
	successCount := 0

	for i := 0; i < numLargeRequests; i++ {
		startTime := time.Now()

		req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		endTime := time.Now()
		requestTime := endTime.Sub(startTime)

		if recorder.Code == http.StatusOK {
			successCount++
			totalTime += requestTime

			// Verify: Large input processing still meets extended performance target
			assert.LessOrEqual(t, requestTime.Milliseconds(), int64(200), "Expected large input processing under 200ms")
		}
	}

	// Verify: All large input requests succeeded
	assert.Equal(t, numLargeRequests, successCount, "Expected all large input requests to succeed")

	// Verify: Average time reasonable for large inputs
	avgTime := totalTime / time.Duration(successCount)
	assert.LessOrEqual(t, avgTime.Milliseconds(), int64(150), "Expected average large input processing under 150ms")

	t.Logf("Large Input Performance: Input Size: %d, Avg Time: %v", inputSize, avgTime)
}

// TestPerformanceMemoryUsage tests memory usage under performance load
func TestPerformanceMemoryUsage(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "memory-test")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create health handler to monitor memory
	healthService := services.NewHealthService(spatialService, integrationConfig)
	healthHandler := handlers.NewHealthMetricsHandler(healthService)
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)
	api.SetupHealthRoutes(router, healthHandler)

	// Get baseline memory usage
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	var baselineHealth map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &baselineHealth)
	require.NoError(t, err)

	// Extract baseline memory
	components := baselineHealth["components"].(map[string]interface{})
	memory := components["memory"].(map[string]interface{})
	baselineMemoryMB := memory["heap_size_mb"].(float64)

	// Execute: Process many requests to test memory usage
	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)
	for i := 0; i < inputSize/50; i++ {
		testInput[i*50] = 1.0
	}

	performanceRequest := PerformanceTestRequest{
		Input: testInput,
	}

	requestBody, err := json.Marshal(performanceRequest)
	require.NoError(t, err)

	const numMemoryRequests = 100
	for i := 0; i < numMemoryRequests; i++ {
		req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected request to succeed during memory test")
	}

	// Check memory usage after processing
	req, err = http.NewRequest("GET", "/api/v1/health", nil)
	require.NoError(t, err)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	var finalHealth map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &finalHealth)
	require.NoError(t, err)

	components = finalHealth["components"].(map[string]interface{})
	memory = components["memory"].(map[string]interface{})
	finalMemoryMB := memory["heap_size_mb"].(float64)

	// Verify: Memory usage reasonable
	memoryIncrease := finalMemoryMB - baselineMemoryMB
	assert.LessOrEqual(t, memoryIncrease, 100.0, "Expected memory increase under 100MB after 100 requests")
	assert.LessOrEqual(t, finalMemoryMB, 500.0, "Expected total memory usage under 500MB")

	t.Logf("Memory Usage: Baseline: %.2f MB, Final: %.2f MB, Increase: %.2f MB",
		baselineMemoryMB, finalMemoryMB, memoryIncrease)
}
