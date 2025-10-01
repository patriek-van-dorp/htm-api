package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/api"
	"github.com/htm-project/neural-api/internal/handlers"
	"github.com/htm-project/neural-api/internal/infrastructure/config"
	"github.com/htm-project/neural-api/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize application
	app, err := initializeApplication(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start server
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Application represents the main application structure.
type Application struct {
	config     *config.Config
	server     *http.Server
	router     *gin.Engine
	shutdownCh chan os.Signal
}

// initializeApplication sets up the application with all dependencies.
func initializeApplication(cfg *config.Config) (*Application, error) {
	// Set Gin mode based on environment
	// Default to debug mode for development
	gin.SetMode(gin.DebugMode)

	// Create Gin engine
	router := gin.New()

	// Initialize metrics collector (simplified implementation)
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
	appRouter := api.NewRouter(
		httpHandler,
		loggingMiddleware,
		errorMiddleware,
		metricsMiddleware,
		corsMiddleware,
	)

	// Setup routes
	if err := appRouter.SetupRoutes(router); err != nil {
		return nil, fmt.Errorf("failed to setup routes: %w", err)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.ReadTimeout * 2,
	}

	// Setup shutdown channel
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	return &Application{
		config:     cfg,
		server:     server,
		router:     router,
		shutdownCh: shutdownCh,
	}, nil
}

// Run starts the HTTP server and handles graceful shutdown.
func (app *Application) Run() error {
	// Start server in a goroutine
	serverErrCh := make(chan error, 1)
	go func() {
		log.Printf("Starting HTM Neural Processing API server on %s", app.config.Server.Address())
		log.Printf("Environment: development")
		log.Printf("Debug mode: true")

		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- fmt.Errorf("server failed to start: %w", err)
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErrCh:
		return err
	case sig := <-app.shutdownCh:
		log.Printf("Received shutdown signal: %v", sig)
		return app.shutdown()
	}
}

// shutdown performs graceful shutdown of the application.
func (app *Application) shutdown() error {
	log.Println("Initiating graceful shutdown...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := app.server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
		return err
	}

	log.Println("Server shutdown completed")
	return nil
}

// SimpleMetricsCollector is a basic implementation of the MetricsCollector interface.
// In a production environment, this would be replaced with a proper metrics system.
type SimpleMetricsCollector struct {
	requestCount    int64
	errorCount      int64
	processingTimes []int64
	responseTimes   []int64
	concurrentCount int
}

// IncrementRequestCount increments the total request counter.
func (smc *SimpleMetricsCollector) IncrementRequestCount() {
	smc.requestCount++
}

// IncrementErrorCount increments the error counter.
func (smc *SimpleMetricsCollector) IncrementErrorCount() {
	smc.errorCount++
}

// RecordProcessingTime records the time taken for processing.
func (smc *SimpleMetricsCollector) RecordProcessingTime(duration int64) {
	smc.processingTimes = append(smc.processingTimes, duration)
}

// RecordResponseTime records the total response time.
func (smc *SimpleMetricsCollector) RecordResponseTime(duration int64) {
	smc.responseTimes = append(smc.responseTimes, duration)
}

// SetConcurrentRequests sets the current number of concurrent requests.
func (smc *SimpleMetricsCollector) SetConcurrentRequests(count int) {
	smc.concurrentCount = count
}

// GetMetrics returns current metrics snapshot.
func (smc *SimpleMetricsCollector) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"request_count":       smc.requestCount,
		"error_count":         smc.errorCount,
		"concurrent_requests": smc.concurrentCount,
		"avg_processing_time": smc.calculateAverage(smc.processingTimes),
		"avg_response_time":   smc.calculateAverage(smc.responseTimes),
	}
}

// Reset resets all metrics.
func (smc *SimpleMetricsCollector) Reset() {
	smc.requestCount = 0
	smc.errorCount = 0
	smc.processingTimes = nil
	smc.responseTimes = nil
	smc.concurrentCount = 0
}

// calculateAverage calculates the average of a slice of int64 values.
func (smc *SimpleMetricsCollector) calculateAverage(values []int64) float64 {
	if len(values) == 0 {
		return 0
	}

	var sum int64
	for _, v := range values {
		sum += v
	}

	return float64(sum) / float64(len(values))
}
