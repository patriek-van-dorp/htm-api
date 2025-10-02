package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/domain/htm"
	"github.com/htm-project/neural-api/internal/handlers"
	"github.com/htm-project/neural-api/internal/infrastructure/config"
	"github.com/htm-project/neural-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SimpleSpatialPoolerStatusResponse represents the actual response structure
type SimpleSpatialPoolerStatusResponse struct {
	Status        string                    `json:"status"`
	Healthy       bool                      `json:"healthy"`
	Instance      map[string]interface{}    `json:"instance"`
	Configuration *htm.SpatialPoolerConfig  `json:"configuration"`
	Metrics       *htm.SpatialPoolerMetrics `json:"metrics"`
	Timestamp     string                    `json:"timestamp"`
}

// TestSpatialPoolerStatusBasic tests the basic spatial pooler status endpoint
func TestSpatialPoolerStatusBasic(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create actual spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "test-instance")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create spatial pooler handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
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
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected 200 OK status")

	// Verify: Parse JSON response
	var response SimpleSpatialPoolerStatusResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Debug: Print response body if there are issues
	if response.Status != "operational" || !response.Healthy {
		t.Logf("Response body: %s", recorder.Body.String())
	}

	// Verify: Basic response structure (accept either operational or degraded initially)
	assert.Contains(t, []string{"operational", "degraded"}, response.Status, "Expected operational or degraded status")
	assert.NotEmpty(t, response.Timestamp, "Expected timestamp")
	assert.NotNil(t, response.Instance, "Expected instance information")
	assert.NotNil(t, response.Configuration, "Expected configuration")
	assert.NotNil(t, response.Metrics, "Expected metrics")

	// Verify: Configuration matches what we set up
	config := response.Configuration
	assert.Equal(t, integrationConfig.Application.SpatialPooler.InputWidth, config.InputWidth, "Expected input width to match")
	assert.Equal(t, integrationConfig.Application.SpatialPooler.ColumnCount, config.ColumnCount, "Expected column count to match")
	assert.Equal(t, integrationConfig.Application.SpatialPooler.SparsityRatio, config.SparsityRatio, "Expected sparsity ratio to match")

	// Verify: Metrics are initialized
	metrics := response.Metrics
	assert.GreaterOrEqual(t, metrics.TotalProcessed, int64(0), "Expected non-negative total processed")
	assert.NotNil(t, metrics.ErrorCounts, "Expected error counts map")
}
