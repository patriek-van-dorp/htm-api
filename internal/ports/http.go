package ports

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/domain/htm"
)

// HTTPHandler defines the interface for HTTP request handlers.
type HTTPHandler interface {
	// ProcessHTMInput handles POST /api/v1/process requests
	ProcessHTMInput(c *gin.Context)

	// HealthCheck handles GET /health requests
	HealthCheck(c *gin.Context)

	// GetMetrics handles GET /metrics requests
	GetMetrics(c *gin.Context)
}

// ProcessHandler defines the interface specifically for HTM processing endpoints.
type ProcessHandler interface {
	// HandleProcess processes an HTM input request
	HandleProcess(ctx context.Context, request *htm.APIRequest) (*htm.APIResponse, error)

	// ValidateRequest validates an incoming API request
	ValidateRequest(request *htm.APIRequest) error

	// CreateSuccessResponse creates a successful response
	CreateSuccessResponse(requestID string, result *htm.ProcessingResult) *htm.APIResponse

	// CreateErrorResponse creates an error response
	CreateErrorResponse(requestID string, err error) *htm.APIResponse
}

// HealthHandler defines the interface for health check endpoints.
type HealthHandler interface {
	// HandleHealthCheck performs health checks and returns status
	HandleHealthCheck(ctx context.Context) (map[string]interface{}, error)

	// CheckDependencies checks the health of all dependencies
	CheckDependencies(ctx context.Context) map[string]bool

	// GetSystemInfo returns basic system information
	GetSystemInfo() map[string]interface{}
}

// MetricsHandler defines the interface for metrics endpoints.
type MetricsHandler interface {
	// HandleMetrics returns current system metrics
	HandleMetrics(ctx context.Context) (map[string]interface{}, error)

	// GetPerformanceMetrics returns performance-related metrics
	GetPerformanceMetrics() map[string]interface{}

	// GetRequestMetrics returns request-related metrics
	GetRequestMetrics() map[string]interface{}

	// GetSystemMetrics returns system-related metrics
	GetSystemMetrics() map[string]interface{}
}

// Middleware defines the interface for HTTP middleware.
type Middleware interface {
	// Apply applies the middleware to a Gin handler
	Apply() gin.HandlerFunc
}

// LoggingMiddleware defines the interface for request logging middleware.
type LoggingMiddleware interface {
	Middleware

	// LogRequest logs incoming requests
	LogRequest(c *gin.Context)

	// LogResponse logs outgoing responses
	LogResponse(c *gin.Context, statusCode int, responseTime int64)
}

// ErrorMiddleware defines the interface for error handling middleware.
type ErrorMiddleware interface {
	Middleware

	// HandleError processes and logs errors
	HandleError(c *gin.Context, err error)

	// HandlePanic recovers from panics and returns appropriate error response
	HandlePanic(c *gin.Context, recovered interface{})
}

// MetricsMiddleware defines the interface for metrics collection middleware.
type MetricsMiddleware interface {
	Middleware

	// RecordRequest records request metrics
	RecordRequest(c *gin.Context)

	// RecordResponse records response metrics
	RecordResponse(c *gin.Context, statusCode int, responseTime int64)
}

// CORSMiddleware defines the interface for CORS handling middleware.
type CORSMiddleware interface {
	Middleware

	// SetCORSHeaders sets appropriate CORS headers
	SetCORSHeaders(c *gin.Context)

	// HandlePreflight handles CORS preflight requests
	HandlePreflight(c *gin.Context)
}

// Router defines the interface for HTTP routing setup.
type Router interface {
	// SetupRoutes configures all application routes
	SetupRoutes(engine *gin.Engine) error

	// RegisterAPIRoutes registers API v1 routes
	RegisterAPIRoutes(group *gin.RouterGroup) error

	// RegisterHealthRoutes registers health check routes
	RegisterHealthRoutes(engine *gin.Engine) error

	// RegisterMetricsRoutes registers metrics routes
	RegisterMetricsRoutes(engine *gin.Engine) error

	// ApplyMiddleware applies middleware to routes
	ApplyMiddleware(engine *gin.Engine) error
}

// RequestBinder defines the interface for binding HTTP requests to domain objects.
type RequestBinder interface {
	// BindHTMRequest binds HTTP request to HTM API request
	BindHTMRequest(c *gin.Context) (*htm.APIRequest, error)

	// ValidateAndBind validates and binds request data
	ValidateAndBind(c *gin.Context, target interface{}) error
}

// ResponseWriter defines the interface for writing HTTP responses.
type ResponseWriter interface {
	// WriteSuccess writes a successful response
	WriteSuccess(c *gin.Context, data interface{}) error

	// WriteError writes an error response
	WriteError(c *gin.Context, err error) error

	// WriteValidationError writes a validation error response
	WriteValidationError(c *gin.Context, validationErr error) error

	// WriteInternalError writes an internal server error response
	WriteInternalError(c *gin.Context, err error) error
}
