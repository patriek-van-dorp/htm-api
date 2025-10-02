package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/handlers"
	"github.com/htm-project/neural-api/internal/infrastructure/config"
	"github.com/htm-project/neural-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ProcessRequest represents a spatial pooler processing request
type ProcessRequest struct {
	EncoderOutput struct {
		Width      int     `json:"width"`
		ActiveBits []int   `json:"active_bits"`
		Sparsity   float64 `json:"sparsity"`
	} `json:"encoder_output"`
	InputWidth      int                    `json:"input_width"`
	InputID         string                 `json:"input_id"`
	LearningEnabled bool                   `json:"learning_enabled"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// TestSpatialPoolerProcessBasic tests the basic spatial pooler processing endpoint
func TestSpatialPoolerProcessBasic(t *testing.T) {
	// Setup: Create integration configuration
	integrationConfig := config.NewDefaultIntegrationConfig()

	// Setup: Create actual spatial pooling service
	spatialService, err := services.NewSpatialPoolingService(integrationConfig.Application.SpatialPooler, "test-process-instance")
	require.NoError(t, err, "Expected spatial pooling service creation to succeed")

	// Setup: Create spatial pooler handler
	spatialHandler := handlers.NewSpatialPoolerHandler(spatialService)

	// Setup: Configure router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create API router (simplified for testing)
	v1 := router.Group("/api/v1")
	spatialGroup := v1.Group("/spatial-pooler")
	spatialGroup.POST("/process", spatialHandler.ProcessSpatialPooler)

	// Prepare: Create test input with proper HTM sparsity (2-5%)
	request := ProcessRequest{
		InputWidth:      1024,
		InputID:         "test-input-1",
		LearningEnabled: true,
		Metadata:        map[string]interface{}{"test": "data"},
	}
	request.EncoderOutput.Width = 1024
	request.EncoderOutput.ActiveBits = []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210} // About 2.1% sparsity
	request.EncoderOutput.Sparsity = 0.021                                                                                                      // Valid HTM sparsity

	requestBody, err := json.Marshal(request)
	require.NoError(t, err)

	// Execute: Make request to process endpoint
	req, err := http.NewRequest("POST", "/api/v1/spatial-pooler/process", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify: Response status (should succeed even if spatial pooler has health issues)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected 200 OK status, response: %s", recorder.Body.String())

	// Verify: Response is valid JSON
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")

	// Verify: Response has expected fields
	assert.Contains(t, response, "input_id", "Expected input_id in response")
	assert.Contains(t, response, "normalized_sdr", "Expected normalized_sdr in response")
	assert.Contains(t, response, "processing_time_ms", "Expected processing_time_ms in response")

	// Verify: Input ID matches
	assert.Equal(t, "test-input-1", response["input_id"], "Expected input ID to match")

	t.Logf("Processing successful: %+v", response)
}
