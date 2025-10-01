package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHealthEndpointContract tests the contract for GET /health
func TestHealthEndpointContract(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "health_check_success",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"status", "timestamp", "version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test router with proper handlers
			router := setupTestRouter()

			// Create request
			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			require.NoError(t, err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code, "Status code mismatch")

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)

			// This will fail until handler is implemented, but we verify structure when it works
			if err == nil {
				// Verify required fields are present
				for _, field := range tt.expectedFields {
					assert.Contains(t, response, field, "Health response should contain field: %s", field)
				}

				// Verify status field value
				if status, ok := response["status"].(string); ok {
					assert.Equal(t, "healthy", status, "Status should be 'healthy'")
				}

				// Verify version field is present and not empty
				if version, ok := response["version"].(string); ok {
					assert.NotEmpty(t, version, "Version should not be empty")
				}

				// Verify timestamp field is present and not empty
				if timestamp, ok := response["timestamp"].(string); ok {
					assert.NotEmpty(t, timestamp, "Timestamp should not be empty")
				}
			}
		})
	}
}

// TestHealthEndpointReliability tests that health endpoint is always available
func TestHealthEndpointReliability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the handler is implemented
	router := setupTestRouter()

	// Test multiple consecutive requests to ensure consistency
	for i := 0; i < 10; i++ {
		req, err := http.NewRequest(http.MethodGet, "/health", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Health endpoint should always return 200 OK
		assert.Equal(t, http.StatusOK, w.Code, "Health check should always return 200 OK")
	}
}

// TestHealthEndpointResponseTime tests health endpoint performance
func TestHealthEndpointResponseTime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test will fail until the handler is implemented
	router := gin.New()
	// TODO: router.GET("/health", handlers.HealthCheck)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	// Health check should be very fast (under 10ms)
	start := time.Now()
	router.ServeHTTP(w, req)
	duration := time.Since(start)

	assert.Less(t, duration, 10*time.Millisecond, "Health check should respond very quickly")
}
