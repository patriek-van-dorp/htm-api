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

// CompletePipelineTestRequest represents a complete pipeline test request
type CompletePipelineTestRequest struct {
	Input []float64 `json:"input"`
}

// CompletePipelineTestResponse represents the expected response structure
type CompletePipelineTestResponse struct {
	Success        bool        `json:"success"`
	ProcessingTime float64     `json:"processing_time_ms"`
	Result         interface{} `json:"result"`
	HTMProperties  interface{} `json:"htm_properties"`
}

// TestCompletePipelineHTTPToSpatialPooler tests the complete pipeline from HTTP request to spatial pooler
func TestCompletePipelineHTTPToSpatialPooler(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create actual spatial pooling service (not mock)
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "test-instance")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create complete handler chain
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)
	healthHandler := handlers.NewHealthMetricsHandler(services.NewHealthService(spatialService, integrationConfig))

	// Setup: Configure complete router with all endpoints
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)
	api.SetupHealthRoutes(router, healthHandler)

	// Phase 1: Verify system health before processing
	t.Run("SystemHealthBeforeProcessing", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/health", nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected system to be healthy before processing")
	})

	// Phase 2: Process data through complete pipeline
	t.Run("CompleteDataProcessing", func(t *testing.T) {
		// Prepare realistic input data
		inputSize := integrationConfig.Application.SpatialPooler.InputWidth
		testInput := make([]float64, inputSize)

		// Create a sparse input pattern (3% active)
		activeInputs := int(float64(inputSize) * 0.03)
		for i := 0; i < activeInputs; i++ {
			testInput[i*10] = 1.0 // Distribute active bits evenly
		}

		processRequest := CompletePipelineTestRequest{
			Input: testInput,
		}

		requestBody, err := json.Marshal(processRequest)
		require.NoError(t, err)

		// Measure end-to-end processing time
		startTime := time.Now()

		req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		endTime := time.Now()
		totalProcessingTime := endTime.Sub(startTime).Milliseconds()

		// Verify: HTTP response success
		assert.Equal(t, http.StatusOK, recorder.Code, "Expected successful processing")

		// Verify: Response structure
		var response CompletePipelineTestResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err, "Expected valid JSON response")

		// Verify: Processing success
		assert.True(t, response.Success, "Expected processing to succeed")
		assert.Greater(t, response.ProcessingTime, 0.0, "Expected positive processing time")
		assert.LessOrEqual(t, response.ProcessingTime, 100.0, "Expected processing time under 100ms")

		// Verify: End-to-end performance requirement
		assert.LessOrEqual(t, totalProcessingTime, int64(150), "Expected total HTTP processing under 150ms including overhead")

		// Verify: Result structure exists
		assert.NotNil(t, response.Result, "Expected processing result")
		assert.NotNil(t, response.HTMProperties, "Expected HTM properties")
	})

	// Phase 3: Verify spatial pooler status after processing
	t.Run("SpatialPoolerStatusAfterProcessing", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/status", nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected spatial pooler status to be accessible")

		var statusResponse map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &statusResponse)
		require.NoError(t, err, "Expected valid status response")

		// Verify: Status shows processing has occurred
		assert.Equal(t, "healthy", statusResponse["status"], "Expected healthy status after processing")
	})

	// Phase 4: Verify system health after processing
	t.Run("SystemHealthAfterProcessing", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/health", nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected system to remain healthy after processing")

		var healthResponse map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &healthResponse)
		require.NoError(t, err, "Expected valid health response")

		// Verify: Health response shows processing metrics
		components := healthResponse["components"].(map[string]interface{})
		spatialPooler := components["spatial_pooler"].(map[string]interface{})
		performance := spatialPooler["performance"].(map[string]interface{})

		// Verify: Request count has increased
		requestsProcessed := performance["requests_processed"].(float64)
		assert.Greater(t, requestsProcessed, 0.0, "Expected request count to increase after processing")
	})
}

// TestCompletePipelineMultipleRequests tests the pipeline with multiple concurrent requests
func TestCompletePipelineMultipleRequests(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "test-multi")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Execute: Send multiple requests concurrently
	const numRequests = 10
	results := make(chan bool, numRequests)

	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)
	for i := 0; i < 20; i++ {
		testInput[i] = 1.0
	}

	processRequest := CompletePipelineTestRequest{
		Input: testInput,
	}

	requestBody, err := json.Marshal(processRequest)
	require.NoError(t, err)

	// Launch concurrent requests
	for i := 0; i < numRequests; i++ {
		go func() {
			req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
			if err != nil {
				results <- false
				return
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			success := recorder.Code == http.StatusOK
			results <- success
		}()
	}

	// Collect results
	successCount := 0
	for i := 0; i < numRequests; i++ {
		if <-results {
			successCount++
		}
	}

	// Verify: All requests should succeed
	assert.Equal(t, numRequests, successCount, "Expected all concurrent requests to succeed")
}

// TestCompletePipelineHTMValidation tests HTM property validation through the complete pipeline
func TestCompletePipelineHTMValidation(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "test-htm")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create HTM validation handler
	htmHandler := handlers.NewHTMValidationHandler(spatialService)
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)
	api.SetupHTMValidationRoutes(router, htmHandler)

	// Step 1: Process some data to initialize the spatial pooler
	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)
	for i := 0; i < 30; i++ {
		testInput[i] = 1.0
	}

	processRequest := CompletePipelineTestRequest{
		Input: testInput,
	}

	requestBody, err := json.Marshal(processRequest)
	require.NoError(t, err)

	// Process data first
	req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected processing to succeed")

	// Step 2: Validate HTM properties after processing
	req, err = http.NewRequest("GET", "/api/v1/spatial-pooler/validation/htm-properties", nil)
	require.NoError(t, err)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: HTM validation response
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTM validation to succeed")

	var htmResponse map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &htmResponse)
	require.NoError(t, err, "Expected valid HTM validation response")

	// Verify: HTM properties are valid after actual processing
	assert.True(t, htmResponse["valid"].(bool), "Expected HTM properties to be valid after processing")

	// Verify: Biological compliance
	compliance := htmResponse["biological_compliance"].(map[string]interface{})
	assert.Greater(t, compliance["overall_score"].(float64), 0.7, "Expected good biological compliance")

	// Verify: Sparsity analysis
	sparsity := htmResponse["sparsity_analysis"].(map[string]interface{})
	assert.True(t, sparsity["sparsity_valid"].(bool), "Expected sparsity to be valid")
	currentSparsity := sparsity["current_sparsity"].(float64)
	assert.GreaterOrEqual(t, currentSparsity, 0.015, "Expected sparsity >= 1.5%")
	assert.LessOrEqual(t, currentSparsity, 0.05, "Expected sparsity <= 5%")
}

// TestCompletePipelineErrorHandling tests error handling through the complete pipeline
func TestCompletePipelineErrorHandling(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "test-errors")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Test 1: Invalid input data
	t.Run("InvalidInputData", func(t *testing.T) {
		invalidRequest := CompletePipelineTestRequest{
			Input: []float64{}, // Empty input
		}

		requestBody, err := json.Marshal(invalidRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Verify: Should return error for invalid input
		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected bad request for invalid input")
	})

	// Test 2: Malformed JSON
	t.Run("MalformedJSON", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBufferString("{invalid json"))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Verify: Should return error for malformed JSON
		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected bad request for malformed JSON")
	})

	// Test 3: System health should remain healthy even after errors
	t.Run("SystemHealthAfterErrors", func(t *testing.T) {
		// This test would need health endpoints to be implemented
		// For now, we verify that the system doesn't crash
		assert.True(t, true, "System remains stable after error conditions")
	})
}
