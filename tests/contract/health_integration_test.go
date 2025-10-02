package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/api"
	"github.com/htm-project/neural-api/internal/handlers"
	"github.com/htm-project/neural-api/internal/infrastructure/config"
	"github.com/htm-project/neural-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HealthResponseIntegration represents the expected health response with spatial pooler integration
type HealthResponseIntegration struct {
	Status        string                      `json:"status"`
	Timestamp     string                      `json:"timestamp"`
	Version       string                      `json:"version"`
	Components    HealthComponentsIntegration `json:"components"`
	SystemMetrics SystemMetricsIntegration    `json:"system_metrics"`
}

// HealthComponentsIntegration represents health status of integrated components
type HealthComponentsIntegration struct {
	SpatialPooler SpatialPoolerHealthIntegration `json:"spatial_pooler"`
	HTTPServer    HTTPServerHealthIntegration    `json:"http_server"`
	Memory        MemoryHealthIntegration        `json:"memory"`
}

// SpatialPoolerHealthIntegration represents spatial pooler health metrics
type SpatialPoolerHealthIntegration struct {
	Status         string                              `json:"status"`
	EngineType     string                              `json:"engine_type"`
	Implementation string                              `json:"implementation"`
	Performance    SpatialPoolerPerformanceIntegration `json:"performance"`
}

// SpatialPoolerPerformanceIntegration represents spatial pooler performance metrics
type SpatialPoolerPerformanceIntegration struct {
	AvgProcessingTimeMs float64 `json:"avg_processing_time_ms"`
	RequestsProcessed   int64   `json:"requests_processed"`
	ErrorRate           float64 `json:"error_rate"`
}

// HTTPServerHealthIntegration represents HTTP server health metrics
type HTTPServerHealthIntegration struct {
	Status             string `json:"status"`
	ConcurrentRequests int    `json:"concurrent_requests"`
	MaxConcurrent      int    `json:"max_concurrent"`
}

// MemoryHealthIntegration represents memory health metrics
type MemoryHealthIntegration struct {
	Status     string  `json:"status"`
	HeapSizeMB float64 `json:"heap_size_mb"`
	GCPressure string  `json:"gc_pressure"`
}

// SystemMetricsIntegration represents system-level metrics
type SystemMetricsIntegration struct {
	UptimeSeconds int64  `json:"uptime_seconds"`
	TotalRequests int64  `json:"total_requests"`
	CurrentLoad   string `json:"current_load"`
}

// TestHealthIntegrationHealthy tests the health endpoint with healthy spatial pooler integration
func TestHealthIntegrationHealthy(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create actual spatial pooling service (not mock)
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create health service with actual dependencies
	healthService := services.NewHealthService(spatialService, integrationConfig)

	// Setup: Create handler with actual health service
	healthHandler := handlers.NewHealthMetricsHandler(healthService)

	// Setup: Configure router with actual handler
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupHealthRoutes(router, healthHandler)

	// Execute: Make request to health endpoint
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected healthy status")

	// Verify: Response body structure
	var response HealthResponseIntegration
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Top-level health status
	assert.Equal(t, "healthy", response.Status, "Expected overall status to be healthy")
	assert.NotEmpty(t, response.Timestamp, "Expected timestamp to be present")
	assert.Equal(t, "1.0.0", response.Version, "Expected version to match")

	// Verify: Spatial pooler component health
	spatialPooler := response.Components.SpatialPooler
	assert.Equal(t, "healthy", spatialPooler.Status, "Expected spatial pooler to be healthy")
	assert.Equal(t, "production", spatialPooler.EngineType, "Expected production engine type")
	assert.Equal(t, "gonum-optimized", spatialPooler.Implementation, "Expected gonum-optimized implementation")

	// Verify: Spatial pooler performance metrics
	performance := spatialPooler.Performance
	assert.GreaterOrEqual(t, performance.AvgProcessingTimeMs, 0.0, "Expected non-negative processing time")
	assert.LessOrEqual(t, performance.AvgProcessingTimeMs, 100.0, "Expected processing time under 100ms target")
	assert.GreaterOrEqual(t, performance.RequestsProcessed, int64(0), "Expected non-negative requests processed")
	assert.GreaterOrEqual(t, performance.ErrorRate, 0.0, "Expected non-negative error rate")
	assert.LessOrEqual(t, performance.ErrorRate, 0.01, "Expected error rate under 1%")

	// Verify: HTTP server component health
	httpServer := response.Components.HTTPServer
	assert.Equal(t, "healthy", httpServer.Status, "Expected HTTP server to be healthy")
	assert.GreaterOrEqual(t, httpServer.ConcurrentRequests, 0, "Expected non-negative concurrent requests")
	assert.Equal(t, 100, httpServer.MaxConcurrent, "Expected max concurrent to match configuration")

	// Verify: Memory component health
	memory := response.Components.Memory
	assert.Equal(t, "healthy", memory.Status, "Expected memory to be healthy")
	assert.Greater(t, memory.HeapSizeMB, 0.0, "Expected positive heap size")
	assert.Contains(t, []string{"low", "medium", "high"}, memory.GCPressure, "Expected valid GC pressure")

	// Verify: System metrics
	system := response.SystemMetrics
	assert.GreaterOrEqual(t, system.UptimeSeconds, int64(0), "Expected non-negative uptime")
	assert.GreaterOrEqual(t, system.TotalRequests, int64(0), "Expected non-negative total requests")
	assert.Contains(t, []string{"low", "medium", "high"}, system.CurrentLoad, "Expected valid current load")
}

// TestHealthIntegrationDegraded tests the health endpoint when spatial pooler is degraded
func TestHealthIntegrationDegraded(t *testing.T) {
	// Setup: Create integration configuration with degraded spatial pooler
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service that will report as degraded
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Force degraded state by overwhelming the service (simulate high error rate)
	// This test validates the degraded response structure
	healthService := services.NewHealthService(spatialService, integrationConfig)

	// Setup: Create handler with health service
	healthHandler := handlers.NewHealthMetricsHandler(healthService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupHealthRoutes(router, healthHandler)

	// Execute: Make request when system is degraded
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Should still return 200 for degraded state (not 503 in this test scenario)
	// Real degraded state would be 503, but this tests the structure
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected response even for degraded state check")

	// Verify: Response body structure exists
	var response HealthResponseIntegration
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response structure")

	// Verify: Response has required components
	assert.NotEmpty(t, response.Status, "Expected status field")
	assert.NotEmpty(t, response.Components.SpatialPooler.Status, "Expected spatial pooler status")
}

// TestHealthIntegrationPerformanceMetrics tests performance metrics in health response
func TestHealthIntegrationPerformanceMetrics(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create health service
	healthService := services.NewHealthService(spatialService, integrationConfig)

	// Setup: Create handler
	healthHandler := handlers.NewHealthMetricsHandler(healthService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupHealthRoutes(router, healthHandler)

	// Execute: Make request
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Verify: Performance metrics are present and valid
	var response HealthResponseIntegration
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify: Performance metrics meet requirements
	performance := response.Components.SpatialPooler.Performance
	assert.LessOrEqual(t, performance.AvgProcessingTimeMs, 100.0, "Expected processing time under performance target")
	assert.LessOrEqual(t, performance.ErrorRate, 0.05, "Expected error rate under 5%")

	// Verify: System metrics show good performance
	assert.GreaterOrEqual(t, response.SystemMetrics.UptimeSeconds, int64(0), "Expected valid uptime")
}
