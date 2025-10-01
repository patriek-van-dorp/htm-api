package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/htm-project/neural-api/internal/domain/htm"
	"github.com/htm-project/neural-api/internal/ports"
)

// HTTPHandlerImpl implements the HTTPHandler interface.
type HTTPHandlerImpl struct {
	processingService ports.ProcessingService
	validationService ports.ValidationService
	metricsCollector  ports.MetricsCollector
	processHandler    ports.ProcessHandler
	healthHandler     ports.HealthHandler
	metricsHandler    ports.MetricsHandler
}

// NewHTTPHandler creates a new HTTP handler.
func NewHTTPHandler(
	processingService ports.ProcessingService,
	validationService ports.ValidationService,
	metricsCollector ports.MetricsCollector,
	processHandler ports.ProcessHandler,
	healthHandler ports.HealthHandler,
	metricsHandler ports.MetricsHandler,
) ports.HTTPHandler {
	return &HTTPHandlerImpl{
		processingService: processingService,
		validationService: validationService,
		metricsCollector:  metricsCollector,
		processHandler:    processHandler,
		healthHandler:     healthHandler,
		metricsHandler:    metricsHandler,
	}
}

// ProcessHTMInput handles POST /api/v1/process requests.
func (h *HTTPHandlerImpl) ProcessHTMInput(c *gin.Context) {
	start := time.Now()
	requestID := uuid.New().String()

	// Add request ID to context for tracking
	c.Set("request_id", requestID)

	defer func() {
		if h.metricsCollector != nil {
			h.metricsCollector.RecordProcessingTime(time.Since(start).Milliseconds())
			h.metricsCollector.IncrementRequestCount()
		}
	}()

	// Parse and validate request
	var apiRequest htm.APIRequest
	if err := c.ShouldBindJSON(&apiRequest); err != nil {
		h.handleError(c, requestID, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	// Set request ID if not provided
	if apiRequest.RequestID == "" {
		apiRequest.RequestID = requestID
	}

	// Check if processHandler is available
	if h.processHandler == nil {
		h.handleError(c, requestID, http.StatusInternalServerError, "Process handler not available", fmt.Errorf("process handler is nil"))
		return
	}

	// Validate request using the process handler
	if err := h.processHandler.ValidateRequest(&apiRequest); err != nil {
		h.handleError(c, requestID, http.StatusBadRequest, "Request validation failed", err)
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Process the request
	response, err := h.processHandler.HandleProcess(ctx, &apiRequest)
	if err != nil {
		h.handleError(c, requestID, http.StatusInternalServerError, "Processing failed", err)
		return
	}

	// Record success metrics
	if h.metricsCollector != nil {
		h.metricsCollector.IncrementRequestCount()
	}

	// Return successful response
	c.JSON(http.StatusOK, response)
}

// HealthCheck handles GET /health requests.
func (h *HTTPHandlerImpl) HealthCheck(c *gin.Context) {
	start := time.Now()

	defer func() {
		if h.metricsCollector != nil {
			h.metricsCollector.RecordProcessingTime(time.Since(start).Milliseconds())
			h.metricsCollector.IncrementRequestCount()
		}
	}()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Get health status
	_, err := h.healthHandler.HandleHealthCheck(ctx)

	// Create simplified contract-compliant response
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	}

	httpStatus := http.StatusOK
	if err != nil {
		response["status"] = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, response)
}

// GetMetrics handles GET /metrics requests.
func (h *HTTPHandlerImpl) GetMetrics(c *gin.Context) {
	start := time.Now()

	defer func() {
		if h.metricsCollector != nil {
			h.metricsCollector.RecordProcessingTime(time.Since(start).Milliseconds())
			h.metricsCollector.IncrementRequestCount()
		}
	}()

	// Get metrics from collector (simplified contract-compliant response)
	var metrics map[string]interface{}
	if h.metricsCollector != nil {
		collectorMetrics := h.metricsCollector.GetMetrics()
		metrics = map[string]interface{}{
			"request_count":       getMetricValue(collectorMetrics, "total_requests", 0),
			"response_times":      []float64{}, // Simplified for contract
			"error_count":         getMetricValue(collectorMetrics, "failed_requests", 0),
			"concurrent_requests": getMetricValue(collectorMetrics, "active_requests", 0),
		}
	} else {
		// Default metrics when collector is not available
		metrics = map[string]interface{}{
			"request_count":       0,
			"response_times":      []float64{},
			"error_count":         0,
			"concurrent_requests": 0,
		}
	}

	c.JSON(http.StatusOK, metrics)
}

// Helper function to get metric values safely
func getMetricValue(metrics map[string]interface{}, key string, defaultValue int) int {
	if value, ok := metrics[key]; ok {
		if intValue, ok := value.(int); ok {
			return intValue
		}
	}
	return defaultValue
}

// handleError handles error responses consistently.
func (h *HTTPHandlerImpl) handleError(c *gin.Context, requestID string, statusCode int, message string, err error) {
	// Record error metrics
	if h.metricsCollector != nil {
		h.metricsCollector.IncrementErrorCount()
	}

	// Create contract-compliant error response based on status code
	if statusCode >= 400 && statusCode < 500 {
		// Client errors (400-499) - validation errors
		errorResponse := map[string]interface{}{
			"error": map[string]interface{}{
				"code":      "VALIDATION_ERROR",
				"message":   message,
				"retryable": false,
			},
			"request_id": requestID,
		}
		c.JSON(statusCode, errorResponse)
	} else {
		// Server errors (500+) - internal errors
		errorResponse := map[string]interface{}{
			"error": map[string]interface{}{
				"code":      "INTERNAL_ERROR",
				"message":   message,
				"retryable": true,
			},
			"request_id": requestID,
		}
		c.JSON(statusCode, errorResponse)
	}
}

// ProcessHandlerImpl implements the ProcessHandler interface.
type ProcessHandlerImpl struct {
	processingService ports.ProcessingService
	validationService ports.ValidationService
	metricsCollector  ports.MetricsCollector
}

// NewProcessHandler creates a new process handler.
func NewProcessHandler(
	processingService ports.ProcessingService,
	validationService ports.ValidationService,
	metricsCollector ports.MetricsCollector,
) ports.ProcessHandler {
	return &ProcessHandlerImpl{
		processingService: processingService,
		validationService: validationService,
		metricsCollector:  metricsCollector,
	}
}

// HandleProcess processes an HTM input request.
func (ph *ProcessHandlerImpl) HandleProcess(ctx context.Context, request *htm.APIRequest) (*htm.APIResponse, error) {
	start := time.Now()

	defer func() {
		if ph.metricsCollector != nil {
			ph.metricsCollector.RecordProcessingTime(time.Since(start).Milliseconds())
		}
	}()

	// Validate the request
	if err := ph.ValidateRequest(request); err != nil {
		return ph.CreateErrorResponse(request.RequestID, err), nil
	}

	// Process the input
	result, err := ph.processingService.ProcessHTMInput(ctx, &request.Input)
	if err != nil {
		return ph.CreateErrorResponse(request.RequestID, err), nil
	}

	// Update result with request ID
	result.ID = request.RequestID

	// Create successful response
	return ph.CreateSuccessResponse(request.RequestID, result), nil
}

// ValidateRequest validates an incoming API request.
func (ph *ProcessHandlerImpl) ValidateRequest(request *htm.APIRequest) error {
	if request == nil {
		return &htm.ValidationError{
			Field:   "request",
			Message: "Request cannot be nil",
		}
	}

	// Use the validation service to properly validate the request
	if ph.validationService == nil {
		return &htm.ValidationError{
			Field:   "validation_service",
			Message: "Validation service not available",
		}
	}

	return ph.validationService.ValidateAPIRequest(request)
}

// CreateSuccessResponse creates a successful response.
func (ph *ProcessHandlerImpl) CreateSuccessResponse(requestID string, result *htm.ProcessingResult) *htm.APIResponse {
	return &htm.APIResponse{
		RequestID:    requestID,
		Result:       result,
		ResponseTime: time.Now(),
	}
}

// CreateErrorResponse creates an error response.
func (ph *ProcessHandlerImpl) CreateErrorResponse(requestID string, err error) *htm.APIResponse {
	apiError := &htm.APIError{
		Code:    "PROCESSING_ERROR",
		Message: err.Error(),
		Details: make(map[string]interface{}),
	}

	// Check if it's a validation error and extract details
	if validationErr, ok := err.(*htm.ValidationError); ok {
		apiError.Code = "VALIDATION_ERROR"
		apiError.Details["field"] = validationErr.Field
		apiError.Details["validation_message"] = validationErr.Message
	}

	return &htm.APIResponse{
		RequestID:    requestID,
		Error:        apiError,
		ResponseTime: time.Now(),
	}
}
