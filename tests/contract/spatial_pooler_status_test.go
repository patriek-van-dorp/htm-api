package contract

import (
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

// SpatialPoolerStatusResponse represents the spatial pooler status response
type SpatialPoolerStatusResponse struct {
	Status        string                    `json:"status"`
	Healthy       bool                      `json:"healthy"`
	Instance      map[string]interface{}    `json:"instance"`
	Configuration *htm.SpatialPoolerConfig  `json:"configuration"` // Use actual HTM config
	Metrics       *htm.SpatialPoolerMetrics `json:"metrics"`       // Use actual HTM metrics
	Timestamp     string                    `json:"timestamp"`
}

// SpatialPoolerEngineStatus represents the engine status
type SpatialPoolerEngineStatus struct {
	Type           string `json:"type"`
	Implementation string `json:"implementation"`
	Version        string `json:"version"`
	State          string `json:"state"`
	Initialized    bool   `json:"initialized"`
}

// SpatialPoolerConfigurationStatus represents the configuration status
type SpatialPoolerConfigurationStatus struct {
	InputWidth       int     `json:"input_width"`
	ColumnCount      int     `json:"column_count"`
	SparsityRatio    float64 `json:"sparsity_ratio"`
	PotentialRadius  int     `json:"potential_radius"`
	GlobalInhibition bool    `json:"global_inhibition"`
	LastUpdated      string  `json:"last_updated"`
	ValidationStatus string  `json:"validation_status"`
}

// SpatialPoolerPerformanceStatus represents the performance status
type SpatialPoolerPerformanceStatus struct {
	AvgProcessingTimeMs  float64 `json:"avg_processing_time_ms"`
	LastProcessingTimeMs float64 `json:"last_processing_time_ms"`
	TotalProcessingCount int64   `json:"total_processing_count"`
	ErrorCount           int64   `json:"error_count"`
	ErrorRate            float64 `json:"error_rate"`
	ThroughputPerSecond  float64 `json:"throughput_per_second"`
	MemoryUsageMB        float64 `json:"memory_usage_mb"`
}

// HTMPropertiesStatus represents HTM biological properties status
type HTMPropertiesStatus struct {
	SparsityValid     bool    `json:"sparsity_valid"`
	CurrentSparsity   float64 `json:"current_sparsity"`
	TargetSparsity    float64 `json:"target_sparsity"`
	OverlapValid      bool    `json:"overlap_valid"`
	LearningActive    bool    `json:"learning_active"`
	InhibitionWorking bool    `json:"inhibition_working"`
	PropertiesHealthy bool    `json:"properties_healthy"`
}

// SpatialPoolerStatistics represents operational statistics
type SpatialPoolerStatistics struct {
	ActiveColumns       int              `json:"active_columns"`
	OverlapDistribution []int            `json:"overlap_distribution"`
	BoostFactorStats    BoostFactorStats `json:"boost_factor_stats"`
	DutyCycleStats      DutyCycleStats   `json:"duty_cycle_stats"`
}

// BoostFactorStats represents boost factor statistics
type BoostFactorStats struct {
	Mean   float64 `json:"mean"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	StdDev float64 `json:"std_dev"`
}

// DutyCycleStats represents duty cycle statistics
type DutyCycleStats struct {
	ActiveMean    float64 `json:"active_mean"`
	OverlapMean   float64 `json:"overlap_mean"`
	MinActiveRate float64 `json:"min_active_rate"`
	MaxActiveRate float64 `json:"max_active_rate"`
}

// TestSpatialPoolerStatusHealthy tests the spatial pooler status endpoint when healthy
func TestSpatialPoolerStatusHealthy(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create actual spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "test-instance")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create spatial pooler handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router with full middleware stack
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create API router (simplified for testing)
	v1 := router.Group("/api/v1")
	spatialGroup := v1.Group("/spatial-pooler")
	spatialGroup.GET("/status", spatialHandler.GetSpatialPoolerStatus)

	// Execute: Make request to status endpoint
	req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/status", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected healthy status response")

	// Verify: Response body structure
	var response SpatialPoolerStatusResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Status and timestamp
	assert.Equal(t, "healthy", response.Status, "Expected healthy status")
	assert.NotEmpty(t, response.Timestamp, "Expected timestamp to be present")

	// Verify: Engine status
	engine := response.Engine
	assert.Equal(t, "production", engine.Type, "Expected production engine type")
	assert.Equal(t, "gonum-optimized", engine.Implementation, "Expected gonum-optimized implementation")
	assert.NotEmpty(t, engine.Version, "Expected version to be present")
	assert.Equal(t, "running", engine.State, "Expected running state")
	assert.True(t, engine.Initialized, "Expected engine to be initialized")

	// Verify: Configuration status
	config := response.Configuration
	assert.Equal(t, integrationConfig.Application.SpatialPooler.InputWidth, config.InputWidth, "Expected input width to match")
	assert.Equal(t, integrationConfig.Application.SpatialPooler.ColumnCount, config.ColumnCount, "Expected column count to match")
	assert.Equal(t, integrationConfig.Application.SpatialPooler.SparsityRatio, config.SparsityRatio, "Expected sparsity ratio to match")
	assert.Equal(t, integrationConfig.Application.SpatialPooler.InhibitionRadius, config.InhibitionRadius, "Expected inhibition radius to match")
	assert.Equal(t, integrationConfig.Application.SpatialPooler.LearningEnabled, config.LearningEnabled, "Expected learning enabled to match")
	assert.NotEmpty(t, config.LastUpdated, "Expected last updated timestamp")
	assert.Equal(t, "valid", config.ValidationStatus, "Expected configuration to be valid")

	// Verify: Performance status
	performance := response.Performance
	assert.GreaterOrEqual(t, performance.AvgProcessingTimeMs, 0.0, "Expected non-negative average processing time")
	assert.LessOrEqual(t, performance.AvgProcessingTimeMs, 100.0, "Expected processing time under performance target")
	assert.GreaterOrEqual(t, performance.TotalProcessingCount, int64(0), "Expected non-negative processing count")
	assert.GreaterOrEqual(t, performance.ErrorCount, int64(0), "Expected non-negative error count")
	assert.GreaterOrEqual(t, performance.ErrorRate, 0.0, "Expected non-negative error rate")
	assert.LessOrEqual(t, performance.ErrorRate, 0.01, "Expected error rate under 1%")
	assert.GreaterOrEqual(t, performance.ThroughputPerSecond, 0.0, "Expected non-negative throughput")
	assert.GreaterOrEqual(t, performance.MemoryUsageMB, 0.0, "Expected non-negative memory usage")

	// Verify: HTM properties
	htmProps := response.HTMProperties
	assert.True(t, htmProps.SparsityValid, "Expected sparsity to be valid")
	assert.GreaterOrEqual(t, htmProps.CurrentSparsity, 0.015, "Expected sparsity >= 1.5%")
	assert.LessOrEqual(t, htmProps.CurrentSparsity, 0.05, "Expected sparsity <= 5%")
	assert.Equal(t, integrationConfig.Application.SpatialPooler.SparsityRatio, htmProps.TargetSparsity, "Expected target sparsity to match config")
	assert.True(t, htmProps.OverlapValid, "Expected overlap to be valid")
	assert.True(t, htmProps.LearningActive, "Expected learning to be active")
	assert.True(t, htmProps.InhibitionWorking, "Expected inhibition to be working")
	assert.True(t, htmProps.PropertiesHealthy, "Expected HTM properties to be healthy")

	// Verify: Statistics
	stats := response.Statistics
	assert.Greater(t, stats.ActiveColumns, 0, "Expected positive active columns")
	expectedActiveColumns := int(float64(integrationConfig.Application.SpatialPooler.ColumnCount) * integrationConfig.Application.SpatialPooler.SparsityRatio)
	assert.InDelta(t, expectedActiveColumns, stats.ActiveColumns, float64(expectedActiveColumns)*0.1, "Expected active columns near target")
	assert.NotEmpty(t, stats.OverlapDistribution, "Expected overlap distribution data")

	// Verify: Boost factor statistics
	boostStats := stats.BoostFactorStats
	assert.GreaterOrEqual(t, boostStats.Mean, 1.0, "Expected boost factor mean >= 1.0")
	assert.GreaterOrEqual(t, boostStats.Min, 1.0, "Expected boost factor min >= 1.0")
	assert.GreaterOrEqual(t, boostStats.Max, boostStats.Mean, "Expected boost factor max >= mean")
	assert.GreaterOrEqual(t, boostStats.StdDev, 0.0, "Expected non-negative std dev")

	// Verify: Duty cycle statistics
	dutyStats := stats.DutyCycleStats
	assert.GreaterOrEqual(t, dutyStats.ActiveMean, 0.0, "Expected non-negative active duty cycle mean")
	assert.GreaterOrEqual(t, dutyStats.OverlapMean, 0.0, "Expected non-negative overlap duty cycle mean")
	assert.GreaterOrEqual(t, dutyStats.MinActiveRate, 0.0, "Expected non-negative min active rate")
	assert.LessOrEqual(t, dutyStats.MaxActiveRate, 1.0, "Expected max active rate <= 1.0")
}

// TestSpatialPoolerStatusUninitialized tests the status endpoint when spatial pooler is uninitialized
func TestSpatialPoolerStatusUninitialized(t *testing.T) {
	// Setup: Create minimal configuration that hasn't been initialized
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial service but don't initialize it properly
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "test-instance-uninit")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create API router (simplified for testing)
	v1 := router.Group("/api/v1")
	spatialGroup := v1.Group("/spatial-pooler")
	spatialGroup.GET("/status", spatialHandler.GetSpatialPoolerStatus)

	// Execute: Make request
	req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/status", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status should still be 200 but status should indicate uninitialized state
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected status response even when uninitialized")

	// Verify: Response body structure
	var response SpatialPoolerStatusResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Should have proper structure even if uninitialized
	assert.NotEmpty(t, response.Status, "Expected status field")
	assert.NotEmpty(t, response.Engine.Type, "Expected engine type")
	assert.NotEmpty(t, response.Configuration.ValidationStatus, "Expected validation status")
}

// TestSpatialPoolerStatusHTMProperties tests HTM biological properties in status
func TestSpatialPoolerStatusHTMProperties(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler)
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create handler
	spatialHandler := handlers.NewSpatialPoolerStatusHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.SetupSpatialPoolerRoutes(router, spatialHandler)

	// Execute: Make request
	req, err := http.NewRequest("GET", "/api/v1/spatial-pooler/status", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Verify: HTM properties meet biological requirements
	var response SpatialPoolerStatusResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)

	htmProps := response.HTMProperties

	// Verify: Sparsity is within biological range (2-5%)
	assert.GreaterOrEqual(t, htmProps.CurrentSparsity, 0.015, "Expected sparsity >= 1.5% (near biological minimum)")
	assert.LessOrEqual(t, htmProps.CurrentSparsity, 0.05, "Expected sparsity <= 5% (biological maximum)")

	// Verify: HTM properties are functioning
	assert.True(t, htmProps.SparsityValid, "Expected sparsity validation to pass")
	assert.True(t, htmProps.OverlapValid, "Expected overlap validation to pass")
	assert.True(t, htmProps.InhibitionWorking, "Expected inhibition mechanism to work")
	assert.True(t, htmProps.PropertiesHealthy, "Expected overall HTM properties to be healthy")
}
