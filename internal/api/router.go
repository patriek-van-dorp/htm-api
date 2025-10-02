package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/ports"
)

// RouterImpl implements the Router interface.
type RouterImpl struct {
	httpHandler          ports.HTTPHandler
	spatialPoolerHandler SpatialPoolerHandler // Add spatial pooler handler
	loggingMiddleware    ports.LoggingMiddleware
	errorMiddleware      ports.ErrorMiddleware
	metricsMiddleware    ports.MetricsMiddleware
	corsMiddleware       ports.CORSMiddleware
}

// SpatialPoolerHandler interface for spatial pooler HTTP operations
type SpatialPoolerHandler interface {
	ProcessSpatialPooler(c *gin.Context)
	GetSpatialPoolerConfig(c *gin.Context)
	UpdateSpatialPoolerConfig(c *gin.Context)
	GetSpatialPoolerMetrics(c *gin.Context)
	ResetSpatialPoolerMetrics(c *gin.Context)
	GetSpatialPoolerStatus(c *gin.Context)
	GetSpatialPoolerHealth(c *gin.Context)
	GetHTMProperties(c *gin.Context)
	ValidateConfigRequest(c *gin.Context)
}

// NewRouter creates a new router.
func NewRouter(
	httpHandler ports.HTTPHandler,
	spatialPoolerHandler SpatialPoolerHandler,
	loggingMiddleware ports.LoggingMiddleware,
	errorMiddleware ports.ErrorMiddleware,
	metricsMiddleware ports.MetricsMiddleware,
	corsMiddleware ports.CORSMiddleware,
) ports.Router {
	return &RouterImpl{
		httpHandler:          httpHandler,
		spatialPoolerHandler: spatialPoolerHandler,
		loggingMiddleware:    loggingMiddleware,
		errorMiddleware:      errorMiddleware,
		metricsMiddleware:    metricsMiddleware,
		corsMiddleware:       corsMiddleware,
	}
}

// NewRouterWithoutSpatialPooler creates a new router without spatial pooler (backward compatibility)
func NewRouterWithoutSpatialPooler(
	httpHandler ports.HTTPHandler,
	loggingMiddleware ports.LoggingMiddleware,
	errorMiddleware ports.ErrorMiddleware,
	metricsMiddleware ports.MetricsMiddleware,
	corsMiddleware ports.CORSMiddleware,
) ports.Router {
	return &RouterImpl{
		httpHandler:          httpHandler,
		spatialPoolerHandler: nil, // No spatial pooler handler
		loggingMiddleware:    loggingMiddleware,
		errorMiddleware:      errorMiddleware,
		metricsMiddleware:    metricsMiddleware,
		corsMiddleware:       corsMiddleware,
	}
}

// SetupRoutes configures all application routes.
func (r *RouterImpl) SetupRoutes(engine *gin.Engine) error {
	// Apply global middleware
	if err := r.ApplyMiddleware(engine); err != nil {
		return err
	}

	// Register health routes (no API versioning)
	if err := r.RegisterHealthRoutes(engine); err != nil {
		return err
	}

	// Register metrics routes (no API versioning)
	if err := r.RegisterMetricsRoutes(engine); err != nil {
		return err
	}

	// Register API v1 routes
	apiV1 := engine.Group("/api/v1")
	if err := r.RegisterAPIRoutes(apiV1); err != nil {
		return err
	}

	// Add a root route for basic info
	engine.GET("/", r.handleRoot)

	return nil
}

// RegisterAPIRoutes registers API v1 routes.
func (r *RouterImpl) RegisterAPIRoutes(group *gin.RouterGroup) error {
	if r.httpHandler == nil {
		return &RouterError{
			Route:   "/api/v1",
			Message: "HTTP handler not available",
		}
	}

	// HTM processing endpoint
	group.POST("/process", r.httpHandler.ProcessHTMInput)

	// Register spatial pooler routes if handler is available
	if r.spatialPoolerHandler != nil {
		if err := r.RegisterSpatialPoolerRoutes(group); err != nil {
			return err
		}
	}

	// Add additional API endpoints here as needed
	// group.GET("/status", r.httpHandler.GetStatus)
	// group.GET("/models", r.httpHandler.ListModels)

	return nil
}

// RegisterSpatialPoolerRoutes registers spatial pooler routes.
func (r *RouterImpl) RegisterSpatialPoolerRoutes(group *gin.RouterGroup) error {
	if r.spatialPoolerHandler == nil {
		return &RouterError{
			Route:   "/api/v1/spatial-pooler",
			Message: "Spatial pooler handler not available",
		}
	}

	// Create spatial pooler route group
	spatialGroup := group.Group("/spatial-pooler")

	// Spatial pooler processing endpoint
	spatialGroup.POST("/process", r.spatialPoolerHandler.ProcessSpatialPooler)

	// Configuration endpoints
	spatialGroup.GET("/config", r.spatialPoolerHandler.GetSpatialPoolerConfig)
	spatialGroup.PUT("/config", r.spatialPoolerHandler.UpdateSpatialPoolerConfig)
	spatialGroup.POST("/config/validate", r.spatialPoolerHandler.ValidateConfigRequest)

	// Validation endpoints
	validationGroup := spatialGroup.Group("/validation")
	validationGroup.GET("/htm-properties", r.spatialPoolerHandler.GetHTMProperties)

	// Metrics endpoints
	spatialGroup.GET("/metrics", r.spatialPoolerHandler.GetSpatialPoolerMetrics)
	spatialGroup.POST("/metrics/reset", r.spatialPoolerHandler.ResetSpatialPoolerMetrics)

	// Status and health endpoints
	spatialGroup.GET("/status", r.spatialPoolerHandler.GetSpatialPoolerStatus)
	spatialGroup.GET("/health", r.spatialPoolerHandler.GetSpatialPoolerHealth)

	return nil
}

// RegisterHealthRoutes registers health check routes.
func (r *RouterImpl) RegisterHealthRoutes(engine *gin.Engine) error {
	if r.httpHandler == nil {
		return &RouterError{
			Route:   "/health",
			Message: "HTTP handler not available",
		}
	}

	// Health check endpoints
	engine.GET("/health", r.httpHandler.HealthCheck)
	engine.GET("/health/ready", r.httpHandler.HealthCheck) // Kubernetes readiness probe
	engine.GET("/health/live", r.httpHandler.HealthCheck)  // Kubernetes liveness probe

	return nil
}

// RegisterMetricsRoutes registers metrics routes.
func (r *RouterImpl) RegisterMetricsRoutes(engine *gin.Engine) error {
	if r.httpHandler == nil {
		return &RouterError{
			Route:   "/metrics",
			Message: "HTTP handler not available",
		}
	}

	// Metrics endpoint
	engine.GET("/metrics", r.httpHandler.GetMetrics)

	return nil
}

// ApplyMiddleware applies middleware to routes.
func (r *RouterImpl) ApplyMiddleware(engine *gin.Engine) error {
	// Recovery middleware (should be first)
	engine.Use(gin.Recovery())

	// CORS middleware
	if r.corsMiddleware != nil {
		engine.Use(r.corsMiddleware.Apply())
	}

	// Logging middleware
	if r.loggingMiddleware != nil {
		engine.Use(r.loggingMiddleware.Apply())
	}

	// Metrics middleware
	if r.metricsMiddleware != nil {
		engine.Use(r.metricsMiddleware.Apply())
	}

	return nil
}

// handleRoot handles the root endpoint.
func (r *RouterImpl) handleRoot(c *gin.Context) {
	endpoints := map[string]string{
		"process": "/api/v1/process",
		"health":  "/health",
		"metrics": "/metrics",
	}

	// Add spatial pooler endpoints if handler is available
	if r.spatialPoolerHandler != nil {
		endpoints["spatial_pooler_process"] = "/api/v1/spatial-pooler/process"
		endpoints["spatial_pooler_config"] = "/api/v1/spatial-pooler/config"
		endpoints["spatial_pooler_metrics"] = "/api/v1/spatial-pooler/metrics"
		endpoints["spatial_pooler_health"] = "/api/v1/spatial-pooler/health"
	}

	c.JSON(http.StatusOK, gin.H{
		"service":       "HTM Neural Processing API",
		"version":       "1.0.0",
		"status":        "running",
		"endpoints":     endpoints,
		"documentation": "https://github.com/htm-project/neural-api",
		"features":      r.getFeatureList(),
	})
}

// getFeatureList returns list of available features
func (r *RouterImpl) getFeatureList() []string {
	features := []string{"htm_processing", "health_monitoring", "metrics"}

	if r.spatialPoolerHandler != nil {
		features = append(features, "spatial_pooler")
	}

	return features
}

// RouterError represents a router configuration error.
type RouterError struct {
	Route   string
	Message string
}

// Error implements the error interface.
func (e *RouterError) Error() string {
	return "Router error for route '" + e.Route + "': " + e.Message
}

// MiddlewareFactory provides methods to create middleware instances.
type MiddlewareFactory struct{}

// NewMiddlewareFactory creates a new middleware factory.
func NewMiddlewareFactory() *MiddlewareFactory {
	return &MiddlewareFactory{}
}

// CreateLoggingMiddleware creates a logging middleware.
func (mf *MiddlewareFactory) CreateLoggingMiddleware() ports.LoggingMiddleware {
	return &LoggingMiddlewareImpl{}
}

// CreateErrorMiddleware creates an error handling middleware.
func (mf *MiddlewareFactory) CreateErrorMiddleware() ports.ErrorMiddleware {
	return &ErrorMiddlewareImpl{}
}

// CreateMetricsMiddleware creates a metrics collection middleware.
func (mf *MiddlewareFactory) CreateMetricsMiddleware(collector ports.MetricsCollector) ports.MetricsMiddleware {
	return &MetricsMiddlewareImpl{collector: collector}
}

// CreateCORSMiddleware creates a CORS handling middleware.
func (mf *MiddlewareFactory) CreateCORSMiddleware() ports.CORSMiddleware {
	return &CORSMiddlewareImpl{}
}

// LoggingMiddlewareImpl implements the LoggingMiddleware interface.
type LoggingMiddlewareImpl struct{}

// Apply applies the logging middleware to a Gin handler.
func (lm *LoggingMiddlewareImpl) Apply() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
	})
}

// LogRequest logs incoming requests.
func (lm *LoggingMiddlewareImpl) LogRequest(c *gin.Context) {
	// Custom request logging logic would go here
	// For now, this is handled by the Gin logger
}

// LogResponse logs outgoing responses.
func (lm *LoggingMiddlewareImpl) LogResponse(c *gin.Context, statusCode int, responseTime int64) {
	// Custom response logging logic would go here
	// For now, this is handled by the Gin logger
}

// ErrorMiddlewareImpl implements the ErrorMiddleware interface.
type ErrorMiddlewareImpl struct{}

// Apply applies the error handling middleware to a Gin handler.
func (em *ErrorMiddlewareImpl) Apply() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				em.HandlePanic(c, recovered)
			}
		}()

		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			em.HandleError(c, err)
		}
	}
}

// HandleError processes and logs errors.
func (em *ErrorMiddlewareImpl) HandleError(c *gin.Context, err error) {
	// Log the error
	// In a real implementation, this would use a proper logger

	// If response hasn't been written yet, write error response
	if !c.Writer.Written() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": err.Error(),
		})
	}
}

// HandlePanic recovers from panics and returns appropriate error response.
func (em *ErrorMiddlewareImpl) HandlePanic(c *gin.Context, recovered interface{}) {
	// Log the panic
	// In a real implementation, this would use a proper logger

	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "Internal server error",
		"details": "An unexpected error occurred",
	})

	c.Abort()
}

// MetricsMiddlewareImpl implements the MetricsMiddleware interface.
type MetricsMiddlewareImpl struct {
	collector ports.MetricsCollector
}

// Apply applies the metrics collection middleware to a Gin handler.
func (mm *MetricsMiddlewareImpl) Apply() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		mm.RecordRequest(c)

		c.Next()

		duration := time.Since(start)
		mm.RecordResponse(c, c.Writer.Status(), duration.Milliseconds())
	}
}

// RecordRequest records request metrics.
func (mm *MetricsMiddlewareImpl) RecordRequest(c *gin.Context) {
	if mm.collector != nil {
		mm.collector.IncrementRequestCount()
	}
}

// RecordResponse records response metrics.
func (mm *MetricsMiddlewareImpl) RecordResponse(c *gin.Context, statusCode int, responseTime int64) {
	if mm.collector != nil {
		mm.collector.RecordProcessingTime(responseTime)
	}
}

// CORSMiddlewareImpl implements the CORSMiddleware interface.
type CORSMiddlewareImpl struct{}

// Apply applies the CORS handling middleware to a Gin handler.
func (cm *CORSMiddlewareImpl) Apply() gin.HandlerFunc {
	return func(c *gin.Context) {
		cm.SetCORSHeaders(c)

		if c.Request.Method == "OPTIONS" {
			cm.HandlePreflight(c)
			return
		}

		c.Next()
	}
}

// SetCORSHeaders sets appropriate CORS headers.
func (cm *CORSMiddlewareImpl) SetCORSHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	c.Header("Access-Control-Max-Age", "86400")
}

// HandlePreflight handles CORS preflight requests.
func (cm *CORSMiddlewareImpl) HandlePreflight(c *gin.Context) {
	c.Status(http.StatusOK)
	c.Abort()
}
