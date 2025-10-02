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
	"github.com/htm-project/neural-api/internal/api"
	"github.com/htm-project/neural-api/internal/handlers"
	"github.com/htm-project/neural-api/internal/infrastructure/config"
	"github.com/htm-project/neural-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ConcurrencyTestRequest represents a concurrent test request
type ConcurrencyTestRequest struct {
	Input []float64 `json:"input"`
	ID    string    `json:"id,omitempty"`
}

// ConcurrencyTestResult represents the result of a concurrent request
type ConcurrencyTestResult struct {
	Success        bool          `json:"success"`
	ProcessingTime time.Duration `json:"processing_time"`
	StatusCode     int           `json:"status_code"`
	RequestID      string        `json:"request_id"`
	Error          string        `json:"error,omitempty"`
}

// TestConcurrentRequestsLimit tests the 100 concurrent requests limit
func TestConcurrentRequestsLimit(t *testing.T) {
	// Setup: Create integration configuration with concurrency limits
	integrationConfig := config.NewDefaultIntegrationConfig()
	integrationConfig.Performance.MaxConcurrentRequests = 100

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "concurrency-test")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router with concurrency limiting middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Add concurrency limiting middleware
	router.Use(handlers.ConcurrencyLimitMiddleware(integrationConfig.Performance.MaxConcurrentRequests))
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Test: Send exactly 100 concurrent requests (should all succeed)
	const exactLimitRequests = 100
	var wg sync.WaitGroup
	results := make(chan ConcurrencyTestResult, exactLimitRequests)

	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)
	for i := 0; i < inputSize/100; i++ {
		testInput[i*100] = 1.0
	}

	// Launch exactly 100 concurrent requests
	for i := 0; i < exactLimitRequests; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			concurrencyRequest := ConcurrencyTestRequest{
				Input: testInput,
				ID:    fmt.Sprintf("req-%d", requestID),
			}

			requestBody, err := json.Marshal(concurrencyRequest)
			if err != nil {
				results <- ConcurrencyTestResult{
					Success:   false,
					RequestID: fmt.Sprintf("req-%d", requestID),
					Error:     err.Error(),
				}
				return
			}

			startTime := time.Now()

			req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
			if err != nil {
				results <- ConcurrencyTestResult{
					Success:   false,
					RequestID: fmt.Sprintf("req-%d", requestID),
					Error:     err.Error(),
				}
				return
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			endTime := time.Now()
			processingTime := endTime.Sub(startTime)

			results <- ConcurrencyTestResult{
				Success:        recorder.Code == http.StatusOK,
				ProcessingTime: processingTime,
				StatusCode:     recorder.Code,
				RequestID:      fmt.Sprintf("req-%d", requestID),
			}
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	close(results)

	// Analyze results
	successCount := 0
	totalProcessingTime := time.Duration(0)
	maxProcessingTime := time.Duration(0)
	var statusCodes []int

	for result := range results {
		statusCodes = append(statusCodes, result.StatusCode)
		if result.Success {
			successCount++
			totalProcessingTime += result.ProcessingTime
			if result.ProcessingTime > maxProcessingTime {
				maxProcessingTime = result.ProcessingTime
			}
		}
	}

	// Verify: All requests within limit should succeed
	assert.GreaterOrEqual(t, successCount, exactLimitRequests*9/10, "Expected at least 90% of concurrent requests to succeed")

	// Verify: Average processing time reasonable under concurrency
	if successCount > 0 {
		avgProcessingTime := totalProcessingTime / time.Duration(successCount)
		assert.LessOrEqual(t, avgProcessingTime.Milliseconds(), int64(200), "Expected average processing time under 200ms with concurrency")
		assert.LessOrEqual(t, maxProcessingTime.Milliseconds(), int64(500), "Expected max processing time under 500ms with concurrency")
	}

	t.Logf("Concurrency Test Results: %d/%d succeeded, Avg Time: %v, Max Time: %v",
		successCount, exactLimitRequests, totalProcessingTime/time.Duration(successCount), maxProcessingTime)
}

// TestConcurrentRequestsExceedLimit tests behavior when exceeding the 100 concurrent requests limit
func TestConcurrentRequestsExceedLimit(t *testing.T) {
	// Setup: Create integration configuration with concurrency limits
	integrationConfig := config.NewDefaultIntegrationConfig()
	integrationConfig.Performance.MaxConcurrentRequests = 50 // Lower limit for testing

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "exceed-limit-test")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router with concurrency limiting
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(handlers.ConcurrencyLimitMiddleware(integrationConfig.Performance.MaxConcurrentRequests))
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Test: Send more requests than the limit
	const exceedLimitRequests = 75 // More than the 50 limit
	var wg sync.WaitGroup
	results := make(chan ConcurrencyTestResult, exceedLimitRequests)

	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)
	for i := 0; i < inputSize/200; i++ {
		testInput[i*200] = 1.0
	}

	// Launch requests that exceed the limit
	for i := 0; i < exceedLimitRequests; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			concurrencyRequest := ConcurrencyTestRequest{
				Input: testInput,
				ID:    fmt.Sprintf("exceed-req-%d", requestID),
			}

			requestBody, err := json.Marshal(concurrencyRequest)
			if err != nil {
				results <- ConcurrencyTestResult{
					Success:   false,
					RequestID: fmt.Sprintf("exceed-req-%d", requestID),
					Error:     err.Error(),
				}
				return
			}

			req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
			if err != nil {
				results <- ConcurrencyTestResult{
					Success:   false,
					RequestID: fmt.Sprintf("exceed-req-%d", requestID),
					Error:     err.Error(),
				}
				return
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			results <- ConcurrencyTestResult{
				Success:    recorder.Code == http.StatusOK,
				StatusCode: recorder.Code,
				RequestID:  fmt.Sprintf("exceed-req-%d", requestID),
			}
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	close(results)

	// Analyze results
	successCount := 0
	tooManyRequestsCount := 0
	otherErrorCount := 0

	for result := range results {
		switch result.StatusCode {
		case http.StatusOK:
			successCount++
		case http.StatusTooManyRequests:
			tooManyRequestsCount++
		default:
			otherErrorCount++
		}
	}

	// Verify: Some requests should be rejected with 429 Too Many Requests
	assert.Greater(t, tooManyRequestsCount, 0, "Expected some requests to be rejected with 429 Too Many Requests")

	// Verify: System should handle the overload gracefully
	assert.LessOrEqual(t, otherErrorCount, exceedLimitRequests/10, "Expected less than 10% other errors")

	// Verify: Some requests should still succeed (those within the limit)
	assert.Greater(t, successCount, 0, "Expected some requests to succeed even when exceeding limit")

	t.Logf("Exceed Limit Results: %d succeeded, %d rejected (429), %d other errors",
		successCount, tooManyRequestsCount, otherErrorCount)
}

// TestConcurrentRequestsDifferentInputs tests concurrent processing with different inputs
func TestConcurrentRequestsDifferentInputs(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "diff-inputs-test")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Test: Send concurrent requests with different input patterns
	const diverseRequests = 20
	var wg sync.WaitGroup
	results := make(chan ConcurrencyTestResult, diverseRequests)

	inputSize := integrationConfig.Application.SpatialPooler.InputWidth

	// Launch concurrent requests with different input patterns
	for i := 0; i < diverseRequests; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			// Create unique input pattern for each request
			testInput := make([]float64, inputSize)

			// Pattern varies based on request ID
			offset := requestID % 50
			for j := 0; j < inputSize/100; j++ {
				testInput[(j*100+offset)%inputSize] = 1.0
			}

			concurrencyRequest := ConcurrencyTestRequest{
				Input: testInput,
				ID:    fmt.Sprintf("diverse-req-%d", requestID),
			}

			requestBody, err := json.Marshal(concurrencyRequest)
			if err != nil {
				results <- ConcurrencyTestResult{
					Success:   false,
					RequestID: fmt.Sprintf("diverse-req-%d", requestID),
					Error:     err.Error(),
				}
				return
			}

			startTime := time.Now()

			req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
			if err != nil {
				results <- ConcurrencyTestResult{
					Success:   false,
					RequestID: fmt.Sprintf("diverse-req-%d", requestID),
					Error:     err.Error(),
				}
				return
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			endTime := time.Now()
			processingTime := endTime.Sub(startTime)

			results <- ConcurrencyTestResult{
				Success:        recorder.Code == http.StatusOK,
				ProcessingTime: processingTime,
				StatusCode:     recorder.Code,
				RequestID:      fmt.Sprintf("diverse-req-%d", requestID),
			}
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	close(results)

	// Analyze results
	successCount := 0
	var processingTimes []time.Duration

	for result := range results {
		if result.Success {
			successCount++
			processingTimes = append(processingTimes, result.ProcessingTime)
		}
	}

	// Verify: All diverse requests should succeed
	assert.Equal(t, diverseRequests, successCount, "Expected all diverse concurrent requests to succeed")

	// Verify: Processing times should be consistent despite different inputs
	if len(processingTimes) > 0 {
		var totalTime time.Duration
		var maxTime time.Duration
		var minTime time.Duration = time.Hour // Start with large value

		for _, pt := range processingTimes {
			totalTime += pt
			if pt > maxTime {
				maxTime = pt
			}
			if pt < minTime {
				minTime = pt
			}
		}

		avgTime := totalTime / time.Duration(len(processingTimes))

		assert.LessOrEqual(t, avgTime.Milliseconds(), int64(150), "Expected average processing time under 150ms for diverse inputs")
		assert.LessOrEqual(t, maxTime.Milliseconds(), int64(300), "Expected max processing time under 300ms for diverse inputs")

		// Verify: Performance variance is reasonable
		variance := maxTime - minTime
		assert.LessOrEqual(t, variance.Milliseconds(), int64(200), "Expected processing time variance under 200ms")

		t.Logf("Diverse Inputs Results: Avg: %v, Min: %v, Max: %v, Variance: %v",
			avgTime, minTime, maxTime, variance)
	}
}

// TestConcurrentRequestsDataIntegrity tests data integrity under concurrent load
func TestConcurrentRequestsDataIntegrity(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "integrity-test")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Test: Send identical requests concurrently and verify deterministic behavior
	const identicalRequests = 30
	var wg sync.WaitGroup
	responses := make(chan map[string]interface{}, identicalRequests)

	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)

	// Create a fixed pattern for deterministic testing
	for i := 0; i < 10; i++ {
		testInput[i*100] = 1.0
	}

	concurrencyRequest := ConcurrencyTestRequest{
		Input: testInput,
		ID:    "deterministic-test",
	}

	requestBody, err := json.Marshal(concurrencyRequest)
	require.NoError(t, err)

	// Launch identical concurrent requests
	for i := 0; i < identicalRequests; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
			if err != nil {
				responses <- nil
				return
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if recorder.Code == http.StatusOK {
				var response map[string]interface{}
				if json.Unmarshal(recorder.Body.Bytes(), &response) == nil {
					responses <- response
				} else {
					responses <- nil
				}
			} else {
				responses <- nil
			}
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	close(responses)

	// Collect and analyze responses
	var validResponses []map[string]interface{}
	for response := range responses {
		if response != nil {
			validResponses = append(validResponses, response)
		}
	}

	// Verify: All requests should succeed
	assert.Equal(t, identicalRequests, len(validResponses), "Expected all identical concurrent requests to succeed")

	// Verify: Responses should be consistent (deterministic behavior)
	if len(validResponses) > 1 {
		firstResponse := validResponses[0]
		firstSparsity := firstResponse["result"].(map[string]interface{})["sparsity"].(float64)

		for i, response := range validResponses[1:] {
			sparsity := response["result"].(map[string]interface{})["sparsity"].(float64)

			// Allow small variance due to floating point precision and concurrent processing
			sparsityDiff := abs(sparsity - firstSparsity)
			assert.LessOrEqual(t, sparsityDiff, 0.001, "Expected consistent sparsity across concurrent identical requests (request %d)", i+1)
		}

		t.Logf("Data Integrity Results: %d consistent responses, Sparsity: %.4f",
			len(validResponses), firstSparsity)
	}
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
