package contract

import (
	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/api"
	"github.com/htm-project/neural-api/internal/handlers"
	"github.com/htm-project/neural-api/internal/services"
)

// SimpleMetricsCollector is a test implementation of MetricsCollector
type SimpleMetricsCollector struct {
	requestCount   int
	errorCount     int
	responseTime   int64
	concurrentReqs int
}

func (s *SimpleMetricsCollector) IncrementRequestCount() {
	s.requestCount++
}

func (s *SimpleMetricsCollector) IncrementErrorCount() {
	s.errorCount++
}

func (s *SimpleMetricsCollector) RecordProcessingTime(duration int64) {
	s.responseTime = duration
}

func (s *SimpleMetricsCollector) RecordResponseTime(duration int64) {
	s.responseTime = duration
}

func (s *SimpleMetricsCollector) SetConcurrentRequests(count int) {
	s.concurrentReqs = count
}

func (s *SimpleMetricsCollector) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"total_requests":           s.requestCount,
		"failed_requests":          s.errorCount,
		"successful_requests":      s.requestCount - s.errorCount,
		"average_response_time_ms": s.responseTime,
		"active_requests":          s.concurrentReqs,
		"requests_per_second":      0,
	}
}

func (s *SimpleMetricsCollector) Reset() {
	s.requestCount = 0
	s.errorCount = 0
	s.responseTime = 0
	s.concurrentReqs = 0
}

// setupTestRouter creates a properly configured router for contract testing
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Initialize test metrics collector
	metricsCollector := &SimpleMetricsCollector{}

	// Initialize services
	matrixProcessor := services.NewMatrixProcessor(metricsCollector)
	validationService := services.NewValidationService(metricsCollector)
	processingService := services.NewProcessingService(matrixProcessor, validationService, metricsCollector)

	// Initialize handlers
	processHandler := handlers.NewProcessHandler(processingService, validationService, metricsCollector)
	healthHandler := handlers.NewHealthHandler(processingService, metricsCollector)
	metricsHandler := handlers.NewMetricsHandler(metricsCollector)
	httpHandler := handlers.NewHTTPHandler(
		processingService,
		validationService,
		metricsCollector,
		processHandler,
		healthHandler,
		metricsHandler,
	)

	// Initialize middleware factory
	middlewareFactory := api.NewMiddlewareFactory()
	loggingMiddleware := middlewareFactory.CreateLoggingMiddleware()
	errorMiddleware := middlewareFactory.CreateErrorMiddleware()
	metricsMiddleware := middlewareFactory.CreateMetricsMiddleware(metricsCollector)
	corsMiddleware := middlewareFactory.CreateCORSMiddleware()

	// Initialize router
	appRouter := api.NewRouterWithoutSpatialPooler(
		httpHandler,
		loggingMiddleware,
		errorMiddleware,
		metricsMiddleware,
		corsMiddleware,
	)

	// Setup routes
	if err := appRouter.SetupRoutes(router); err != nil {
		panic("Failed to setup routes: " + err.Error())
	}

	return router
}
