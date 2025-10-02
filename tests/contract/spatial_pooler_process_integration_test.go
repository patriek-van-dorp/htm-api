package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/api"
	"github.com/htm-project/neural-api/internal/domain/htm"
	"github.com/htm-project/neural-api/internal/handlers"
	"github.com/htm-project/neural-api/internal/infrastructure/config"
	"github.com/htm-project/neural-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SpatialPoolerProcessRequest represents the process request
type SpatialPoolerProcessRequest struct {
	Input          []float64                `json:"input" validate:"required"`
	Configuration  *htm.SpatialPoolerConfig `json:"configuration,omitempty"`
	ProcessingMode string                   `json:"processing_mode,omitempty"`
	ReturnDetails  bool                     `json:"return_details,omitempty"`
}

// SpatialPoolerProcessResponse represents the process response
type SpatialPoolerProcessResponse struct {
	Success        bool                       `json:"success"`
	ProcessingTime float64                    `json:"processing_time_ms"`
	Timestamp      string                     `json:"timestamp"`
	Result         SpatialPoolerProcessResult `json:"result"`
	HTMProperties  HTMProcessingProperties    `json:"htm_properties"`
	Metrics        ProcessingMetrics          `json:"metrics"`
}

// SpatialPoolerProcessResult represents the processing result
type SpatialPoolerProcessResult struct {
	ActiveColumns []int     `json:"active_columns"`
	SDR           []int     `json:"sdr"`
	Sparsity      float64   `json:"sparsity"`
	OverlapScores []float64 `json:"overlap_scores,omitempty"`
	BoostFactors  []float64 `json:"boost_factors,omitempty"`
}

// HTMProcessingProperties represents HTM properties for this processing
type HTMProcessingProperties struct {
	SparsityValid       bool    `json:"sparsity_valid"`
	SparsityAchieved    float64 `json:"sparsity_achieved"`
	TargetSparsity      float64 `json:"target_sparsity"`
	OverlapPatternValid bool    `json:"overlap_pattern_valid"`
	InhibitionApplied   bool    `json:"inhibition_applied"`
	LearningOccurred    bool    `json:"learning_occurred"`
	BiologicallyValid   bool    `json:"biologically_valid"`
}

// ProcessingMetrics represents performance metrics for this processing
type ProcessingMetrics struct {
	InputSize          int     `json:"input_size"`
	ActiveColumnCount  int     `json:"active_column_count"`
	OverlapComputeTime float64 `json:"overlap_compute_time_ms"`
	InhibitionTime     float64 `json:"inhibition_time_ms"`
	LearningTime       float64 `json:"learning_time_ms"`
	MemoryUsedMB       float64 `json:"memory_used_mb"`
}

// TestSpatialPoolerProcessIntegrationBasic tests basic spatial pooler processing with actual engine
func TestSpatialPoolerProcessIntegrationBasic(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create actual spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create spatial pooler handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Prepare: Create test input data
	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)

	// Create a sparse pattern (5% active bits)
	activeInputBits := int(float64(inputSize) * 0.05)
	for i := 0; i < activeInputBits; i++ {
		testInput[i] = 1.0
	}

	processRequest := SpatialPoolerProcessRequest{
		Input:          testInput,
		ProcessingMode: "normal",
		ReturnDetails:  true,
	}

	// Execute: Make process request
	requestBody, err := json.Marshal(processRequest)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected successful processing")

	// Verify: Response body structure
	var response SpatialPoolerProcessResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Processing success
	assert.True(t, response.Success, "Expected processing to succeed")
	assert.Greater(t, response.ProcessingTime, 0.0, "Expected positive processing time")
	assert.LessOrEqual(t, response.ProcessingTime, 100.0, "Expected processing time under 100ms target")
	assert.NotEmpty(t, response.Timestamp, "Expected timestamp to be present")

	// Verify: Processing result
	result := response.Result
	assert.NotEmpty(t, result.ActiveColumns, "Expected active columns to be returned")
	assert.NotEmpty(t, result.SDR, "Expected SDR to be returned")
	assert.Greater(t, result.Sparsity, 0.0, "Expected positive sparsity")
	assert.LessOrEqual(t, result.Sparsity, 0.1, "Expected sparsity under 10%")

	// Verify: Active columns count matches sparsity
	expectedActiveColumns := int(float64(integrationConfig.Application.SpatialPooler.ColumnCount) * integrationConfig.Application.SpatialPooler.SparsityRatio)
	actualActiveColumns := len(result.ActiveColumns)
	tolerance := float64(expectedActiveColumns) * 0.1 // 10% tolerance
	assert.InDelta(t, expectedActiveColumns, actualActiveColumns, tolerance, "Expected active columns count near target sparsity")

	// Verify: SDR matches active columns
	assert.Equal(t, len(result.ActiveColumns), len(result.SDR), "Expected SDR length to match active columns count")

	// Verify: HTM properties
	htmProps := response.HTMProperties
	assert.True(t, htmProps.SparsityValid, "Expected sparsity to be valid")
	assert.GreaterOrEqual(t, htmProps.SparsityAchieved, 0.015, "Expected sparsity >= 1.5%")
	assert.LessOrEqual(t, htmProps.SparsityAchieved, 0.05, "Expected sparsity <= 5%")
	assert.Equal(t, integrationConfig.Application.SpatialPooler.SparsityRatio, htmProps.TargetSparsity, "Expected target sparsity to match config")
	assert.True(t, htmProps.OverlapPatternValid, "Expected overlap pattern to be valid")
	assert.True(t, htmProps.InhibitionApplied, "Expected inhibition to be applied")
	assert.True(t, htmProps.LearningOccurred, "Expected learning to occur")
	assert.True(t, htmProps.BiologicallyValid, "Expected biologically valid processing")

	// Verify: Processing metrics
	metrics := response.Metrics
	assert.Equal(t, inputSize, metrics.InputSize, "Expected input size to match")
	assert.Equal(t, len(result.ActiveColumns), metrics.ActiveColumnCount, "Expected active column count to match")
	assert.Greater(t, metrics.OverlapComputeTime, 0.0, "Expected positive overlap compute time")
	assert.Greater(t, metrics.InhibitionTime, 0.0, "Expected positive inhibition time")
	assert.GreaterOrEqual(t, metrics.LearningTime, 0.0, "Expected non-negative learning time")
	assert.Greater(t, metrics.MemoryUsedMB, 0.0, "Expected positive memory usage")

	// Verify: Performance requirements
	totalTime := metrics.OverlapComputeTime + metrics.InhibitionTime + metrics.LearningTime
	assert.LessOrEqual(t, totalTime, 100.0, "Expected total processing time under 100ms")
}

// TestSpatialPoolerProcessIntegrationDeterministic tests deterministic behavior
func TestSpatialPoolerProcessIntegrationDeterministic(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Prepare: Create test input data
	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)
	for i := 0; i < 10; i++ {
		testInput[i] = 1.0 // Set first 10 bits active
	}

	processRequest := SpatialPoolerProcessRequest{
		Input:          testInput,
		ProcessingMode: "deterministic",
		ReturnDetails:  true,
	}

	// Execute: First processing
	requestBody, err := json.Marshal(processRequest)
	require.NoError(t, err)

	req1, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req1.Header.Set("Content-Type", "application/json")

	recorder1 := httptest.NewRecorder()
	router.ServeHTTP(recorder1, req1)

	// Execute: Second processing with identical input
	req2, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")

	recorder2 := httptest.NewRecorder()
	router.ServeHTTP(recorder2, req2)

	// Verify: Both requests succeed
	assert.Equal(t, http.StatusOK, recorder1.Code, "Expected first request to succeed")
	assert.Equal(t, http.StatusOK, recorder2.Code, "Expected second request to succeed")

	// Verify: Parse responses
	var response1, response2 SpatialPoolerProcessResponse
	err = json.Unmarshal(recorder1.Body.Bytes(), &response1)
	require.NoError(t, err)
	err = json.Unmarshal(recorder2.Body.Bytes(), &response2)
	require.NoError(t, err)

	// Verify: Deterministic behavior - identical inputs should produce similar sparsity
	sparsityDiff := abs(response1.Result.Sparsity - response2.Result.Sparsity)
	assert.LessOrEqual(t, sparsityDiff, 0.01, "Expected sparsity to be consistent across identical inputs")

	// Verify: HTM properties remain consistent
	assert.Equal(t, response1.HTMProperties.SparsityValid, response2.HTMProperties.SparsityValid, "Expected sparsity validity to be consistent")
	assert.Equal(t, response1.HTMProperties.BiologicallyValid, response2.HTMProperties.BiologicallyValid, "Expected biological validity to be consistent")
}

// TestSpatialPoolerProcessIntegrationLargeInput tests processing with large input (near limit)
func TestSpatialPoolerProcessIntegrationLargeInput(t *testing.T) {
	// Setup: Create integration configuration with large input
	integrationConfig := config.NewDefaultIntegrationConfig()
	integrationConfig.Application.SpatialPooler.InputWidth = 8192 // Large input size

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Prepare: Create large test input data
	inputSize := integrationConfig.Application.SpatialPooler.InputWidth
	testInput := make([]float64, inputSize)

	// Create a pattern with 3% active bits
	activeInputBits := int(float64(inputSize) * 0.03)
	for i := 0; i < activeInputBits; i++ {
		testInput[i] = 1.0
	}

	processRequest := SpatialPoolerProcessRequest{
		Input:          testInput,
		ProcessingMode: "performance",
		ReturnDetails:  true,
	}

	// Execute: Make process request
	requestBody, err := json.Marshal(processRequest)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected successful processing of large input")

	// Verify: Response body
	var response SpatialPoolerProcessResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Processing performance with large input
	assert.True(t, response.Success, "Expected large input processing to succeed")
	assert.LessOrEqual(t, response.ProcessingTime, 100.0, "Expected large input processing under 100ms")

	// Verify: HTM properties maintained with large input
	htmProps := response.HTMProperties
	assert.True(t, htmProps.SparsityValid, "Expected sparsity to remain valid with large input")
	assert.True(t, htmProps.BiologicallyValid, "Expected biological validity with large input")

	// Verify: Memory usage reasonable
	assert.LessOrEqual(t, response.Metrics.MemoryUsedMB, 100.0, "Expected memory usage under 100MB for large input")
}

// TestSpatialPoolerProcessIntegrationInvalidInput tests error handling with invalid input
func TestSpatialPoolerProcessIntegrationInvalidInput(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Test: Empty input
	processRequest := SpatialPoolerProcessRequest{
		Input:          []float64{},
		ProcessingMode: "normal",
		ReturnDetails:  true,
	}

	requestBody, err := json.Marshal(processRequest)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Should return error for empty input
	assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected bad request for empty input")
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
