package integration

import (
	"time"

	"github.com/gin-gonic/gin"
)

// generateTestMatrix creates a test matrix with specified dimensions
func generateTestMatrix(rows, cols int) [][]float64 {
	matrix := make([][]float64, rows)
	for i := range matrix {
		matrix[i] = make([]float64, cols)
		for j := range matrix[i] {
			matrix[i][j] = float64(i*cols + j + 1)
		}
	}
	return matrix
}

// calculateAverageResponseTime calculates the average of response times
func calculateAverageResponseTime(times []time.Duration) time.Duration {
	if len(times) == 0 {
		return 0
	}

	var total time.Duration
	for _, t := range times {
		total += t
	}
	return total / time.Duration(len(times))
}

// setupTestRouter creates a test router with the API routes
func setupTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Add basic health route for testing
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Add basic metrics route for testing
	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"uptime": "5m",
			"status": "operational",
		})
	})

	// Add API routes for testing (these will be implemented later)
	api := router.Group("/api/v1")
	{
		api.POST("/process", func(c *gin.Context) {
			// Simple test implementation
			c.JSON(200, gin.H{
				"status":     "processed",
				"request_id": "test-123",
			})
		})
	}

	return router
}
