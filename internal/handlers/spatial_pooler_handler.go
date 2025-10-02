package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/htm-project/neural-api/internal/domain/htm"
	"github.com/htm-project/neural-api/internal/ports"
)

// SpatialPoolerHandler handles HTTP requests for spatial pooling operations
type SpatialPoolerHandler struct {
	spatialPoolingService ports.SpatialPoolingService
}

// NewSpatialPoolerHandler creates a new spatial pooler HTTP handler
func NewSpatialPoolerHandler(spatialPoolingService ports.SpatialPoolingService) *SpatialPoolerHandler {
	return &SpatialPoolerHandler{
		spatialPoolingService: spatialPoolingService,
	}
}

// ProcessSpatialPooler handles POST /api/v1/spatial-pooler/process requests
func (h *SpatialPoolerHandler) ProcessSpatialPooler(c *gin.Context) {
	var request SpatialPoolerProcessRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validateProcessRequest(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Request validation failed",
			"details": err.Error(),
		})
		return
	}

	// Convert request to domain input
	poolingInput := &htm.PoolingInput{
		EncoderOutput: htm.EncoderOutput{
			Width:      request.EncoderOutput.Width,
			ActiveBits: request.EncoderOutput.ActiveBits,
			Sparsity:   request.EncoderOutput.Sparsity,
		},
		InputWidth:      request.InputWidth,
		InputID:         request.InputID,
		LearningEnabled: request.LearningEnabled,
		Metadata:        request.Metadata,
	}

	// Process with spatial pooling service
	result, err := h.spatialPoolingService.ProcessSpatialPooling(c.Request.Context(), poolingInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Spatial pooling processing failed",
			"details": err.Error(),
		})
		return
	}

	// Convert result to response
	response := SpatialPoolerProcessResponse{
		NormalizedSDR: SDRResponse{
			Width:      result.NormalizedSDR.Width,
			ActiveBits: result.NormalizedSDR.ActiveBits,
			Sparsity:   result.NormalizedSDR.Sparsity,
		},
		InputID:          result.InputID,
		ProcessingTimeMs: result.ProcessingTime,
		ActiveColumns:    result.ActiveColumns,
		AvgOverlap:       result.AvgOverlap,
		SparsityLevel:    result.SparsityLevel,
		LearningOccurred: result.LearningOccurred,
		BoostingApplied:  result.BoostingApplied,
	}

	c.JSON(http.StatusOK, response)
}

// GetSpatialPoolerConfig handles GET /api/v1/spatial-pooler/config requests
func (h *SpatialPoolerHandler) GetSpatialPoolerConfig(c *gin.Context) {
	config, err := h.spatialPoolingService.GetConfiguration(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get configuration",
			"details": err.Error(),
		})
		return
	}

	// Convert to response format
	response := SpatialPoolerConfigResponse{
		InputWidth:          config.InputWidth,
		ColumnCount:         config.ColumnCount,
		SparsityRatio:       config.SparsityRatio,
		Mode:                string(config.Mode),
		LearningEnabled:     config.LearningEnabled,
		LearningRate:        config.LearningRate,
		MaxBoost:            config.MaxBoost,
		BoostStrength:       config.BoostStrength,
		InhibitionRadius:    config.InhibitionRadius,
		LocalAreaDensity:    config.LocalAreaDensity,
		MinOverlapThreshold: config.MinOverlapThreshold,
		MaxProcessingTimeMs: config.MaxProcessingTimeMs,
		SemanticThresholds: SemanticThresholdsResponse{
			SimilarInputMinOverlap:   config.SemanticThresholds.SimilarInputMinOverlap,
			DifferentInputMaxOverlap: config.SemanticThresholds.DifferentInputMaxOverlap,
		},
	}

	c.JSON(http.StatusOK, response)
}

// UpdateSpatialPoolerConfig handles PUT /api/v1/spatial-pooler/config requests
func (h *SpatialPoolerHandler) UpdateSpatialPoolerConfig(c *gin.Context) {
	var request SpatialPoolerConfigUpdateRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Convert request to domain config
	config := &htm.SpatialPoolerConfig{
		InputWidth:          request.InputWidth,
		ColumnCount:         request.ColumnCount,
		SparsityRatio:       request.SparsityRatio,
		LearningEnabled:     request.LearningEnabled,
		LearningRate:        request.LearningRate,
		MaxBoost:            request.MaxBoost,
		BoostStrength:       request.BoostStrength,
		InhibitionRadius:    request.InhibitionRadius,
		LocalAreaDensity:    request.LocalAreaDensity,
		MinOverlapThreshold: request.MinOverlapThreshold,
		MaxProcessingTimeMs: request.MaxProcessingTimeMs,
	}

	// Parse mode
	mode, err := htm.ParseSpatialPoolerMode(request.Mode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid mode value",
			"details": err.Error(),
		})
		return
	}
	config.Mode = mode

	// Set semantic thresholds
	config.SemanticThresholds = htm.SemanticThresholds{
		SimilarInputMinOverlap:   request.SemanticThresholds.SimilarInputMinOverlap,
		DifferentInputMaxOverlap: request.SemanticThresholds.DifferentInputMaxOverlap,
	}

	// Update configuration
	if err := h.spatialPoolingService.UpdateConfiguration(c.Request.Context(), config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Configuration update failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration updated successfully",
	})
}

// GetSpatialPoolerMetrics handles GET /api/v1/spatial-pooler/metrics requests
func (h *SpatialPoolerHandler) GetSpatialPoolerMetrics(c *gin.Context) {
	metrics, err := h.spatialPoolingService.GetMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get metrics",
			"details": err.Error(),
		})
		return
	}

	// Convert to response format
	response := SpatialPoolerMetricsResponse{
		TotalProcessed:           metrics.TotalProcessed,
		AverageProcessingTimeMs:  metrics.AverageProcessingTime,
		LearningIterations:       metrics.LearningIterations,
		ColumnUsageDistribution:  metrics.ColumnUsageDistribution,
		AverageSparsity:          metrics.AverageSparsity,
		OverlapScoreDistribution: metrics.OverlapScoreDistribution,
		BoostingEvents:           metrics.BoostingEvents,
		ErrorCounts:              make(map[string]int64),
	}

	// Convert error counts
	for errorType, count := range metrics.ErrorCounts {
		response.ErrorCounts[string(errorType)] = count
	}

	c.JSON(http.StatusOK, response)
}

// ResetSpatialPoolerMetrics handles POST /api/v1/spatial-pooler/metrics/reset requests
func (h *SpatialPoolerHandler) ResetSpatialPoolerMetrics(c *gin.Context) {
	if err := h.spatialPoolingService.ResetMetrics(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reset metrics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Metrics reset successfully",
	})
}

// GetSpatialPoolerHealth handles GET /api/v1/spatial-pooler/health requests
func (h *SpatialPoolerHandler) GetSpatialPoolerHealth(c *gin.Context) {
	if err := h.spatialPoolingService.HealthCheck(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	// Get instance info
	info := h.spatialPoolingService.GetInstanceInfo(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"info":   info,
	})
}

// GetSpatialPoolerStatus handles GET /api/v1/spatial-pooler/status requests
func (h *SpatialPoolerHandler) GetSpatialPoolerStatus(c *gin.Context) {
	// Get instance information
	info := h.spatialPoolingService.GetInstanceInfo(c.Request.Context())

	// Get current configuration
	config, err := h.spatialPoolingService.GetConfiguration(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get spatial pooler configuration",
			"details": err.Error(),
		})
		return
	}

	// Get current metrics
	metrics, err := h.spatialPoolingService.GetMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get spatial pooler metrics",
			"details": err.Error(),
		})
		return
	}

	// Check health
	isHealthy := true
	var healthError string
	if err := h.spatialPoolingService.HealthCheck(c.Request.Context()); err != nil {
		isHealthy = false
		healthError = err.Error()
	}

	// Compile status response
	status := gin.H{
		"status":        "operational",
		"healthy":       isHealthy,
		"instance":      info,
		"configuration": config,
		"metrics":       metrics,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
	}

	if !isHealthy {
		status["health_error"] = healthError
		status["status"] = "degraded"
	}

	c.JSON(http.StatusOK, status)
}

// ValidateConfigRequest handles POST /api/v1/spatial-pooler/config/validate requests
func (h *SpatialPoolerHandler) ValidateConfigRequest(c *gin.Context) {
	var request SpatialPoolerConfigUpdateRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Convert request to domain config
	config := &htm.SpatialPoolerConfig{
		InputWidth:          request.InputWidth,
		ColumnCount:         request.ColumnCount,
		SparsityRatio:       request.SparsityRatio,
		LearningEnabled:     request.LearningEnabled,
		LearningRate:        request.LearningRate,
		MaxBoost:            request.MaxBoost,
		BoostStrength:       request.BoostStrength,
		InhibitionRadius:    request.InhibitionRadius,
		LocalAreaDensity:    request.LocalAreaDensity,
		MinOverlapThreshold: request.MinOverlapThreshold,
		MaxProcessingTimeMs: request.MaxProcessingTimeMs,
	}

	// Parse mode
	mode, err := htm.ParseSpatialPoolerMode(request.Mode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid":   false,
			"error":   "Invalid mode value",
			"details": err.Error(),
		})
		return
	}
	config.Mode = mode

	// Set semantic thresholds
	config.SemanticThresholds = htm.SemanticThresholds{
		SimilarInputMinOverlap:   request.SemanticThresholds.SimilarInputMinOverlap,
		DifferentInputMaxOverlap: request.SemanticThresholds.DifferentInputMaxOverlap,
	}

	// Validate configuration
	if err := h.spatialPoolingService.ValidateConfiguration(c.Request.Context(), config); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"valid":   false,
			"error":   "Configuration validation failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"message": "Configuration is valid",
	})
}

// GetHTMProperties handles GET /api/v1/spatial-pooler/validation/htm-properties requests
func (h *SpatialPoolerHandler) GetHTMProperties(c *gin.Context) {
	// Get current configuration
	config, err := h.spatialPoolingService.GetConfiguration(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get spatial pooler configuration",
			"details": err.Error(),
		})
		return
	}

	// Get current metrics for runtime properties
	metrics, err := h.spatialPoolingService.GetMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get spatial pooler metrics",
			"details": err.Error(),
		})
		return
	}

	// Calculate HTM compliance properties
	properties := gin.H{
		"htm_compliance": gin.H{
			"biological_constraints": gin.H{
				"sparsity_percentage":   config.SparsityRatio * 100,
				"target_sparsity_range": []float64{2.0, 5.0}, // 2-5% as per HTM theory
				"sparsity_compliant":    config.SparsityRatio >= 0.02 && config.SparsityRatio <= 0.05,
				"overlap_threshold":     config.MinOverlapThreshold,
				"overlap_compliant":     config.MinOverlapThreshold > 0,
			},
			"learning_properties": gin.H{
				"learning_enabled":   config.LearningEnabled,
				"adaptation_enabled": config.LearningEnabled,
				"learning_rate":      config.LearningRate,
				"boost_strength":     config.BoostStrength,
				"max_boost":          config.MaxBoost,
				"learning_compliant": config.LearningEnabled && config.LearningRate > 0,
			},
			"topology_properties": gin.H{
				"column_count":       config.ColumnCount,
				"input_width":        config.InputWidth,
				"inhibition_radius":  config.InhibitionRadius,
				"local_area_density": config.LocalAreaDensity,
				"topology_compliant": config.ColumnCount > 0 && config.InputWidth > 0,
			},
		},
		"runtime_metrics": gin.H{
			"current_sparsity":        metrics.AverageSparsity * 100, // Convert to percentage
			"total_processed":         metrics.TotalProcessed,
			"average_processing_time": metrics.AverageProcessingTime,
			"learning_iterations":     metrics.LearningIterations,
			"boosting_events":         metrics.BoostingEvents,
		},
		"validation_status": gin.H{
			"overall_compliant": h.validateOverallHTMCompliance(config, metrics),
			"warnings":          h.generateHTMWarnings(config, metrics),
			"recommendations":   h.generateHTMRecommendations(config, metrics),
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, properties)
}

// validateOverallHTMCompliance checks if the spatial pooler meets HTM biological principles
func (h *SpatialPoolerHandler) validateOverallHTMCompliance(config *htm.SpatialPoolerConfig, metrics *htm.SpatialPoolerMetrics) bool {
	// Check sparsity compliance (2-5%)
	sparsityCompliant := config.SparsityRatio >= 0.02 && config.SparsityRatio <= 0.05

	// Check if learning is properly configured
	learningCompliant := config.LearningEnabled && config.LearningRate > 0

	// Check topology is valid
	topologyCompliant := config.ColumnCount > 0 && config.InputWidth > 0

	// Check overlap threshold
	overlapCompliant := config.MinOverlapThreshold > 0

	return sparsityCompliant && learningCompliant && topologyCompliant && overlapCompliant
}

// generateHTMWarnings generates warnings for HTM compliance issues
func (h *SpatialPoolerHandler) generateHTMWarnings(config *htm.SpatialPoolerConfig, metrics *htm.SpatialPoolerMetrics) []string {
	warnings := []string{}

	// Check sparsity
	if config.SparsityRatio < 0.02 {
		warnings = append(warnings, "Sparsity ratio below HTM recommended minimum of 2%")
	} else if config.SparsityRatio > 0.05 {
		warnings = append(warnings, "Sparsity ratio above HTM recommended maximum of 5%")
	}

	// Check learning
	if !config.LearningEnabled {
		warnings = append(warnings, "Learning is disabled - HTM requires adaptation")
	} else if config.LearningRate <= 0 {
		warnings = append(warnings, "Learning rate is zero or negative")
	}

	// Check performance
	if metrics.AverageProcessingTime > 10 {
		warnings = append(warnings, "Average processing time exceeds 10ms requirement")
	}

	// Check error rates
	if len(metrics.ErrorCounts) > 0 {
		warnings = append(warnings, "Processing errors detected - check spatial pooler stability")
	}

	return warnings
}

// generateHTMRecommendations generates recommendations for HTM optimization
func (h *SpatialPoolerHandler) generateHTMRecommendations(config *htm.SpatialPoolerConfig, metrics *htm.SpatialPoolerMetrics) []string {
	recommendations := []string{}

	// Sparsity recommendations
	if config.SparsityRatio < 0.02 {
		recommendations = append(recommendations, "Increase sparsity ratio to 2-5% for better HTM compliance")
	} else if config.SparsityRatio > 0.05 {
		recommendations = append(recommendations, "Decrease sparsity ratio to 2-5% for optimal HTM behavior")
	}

	// Learning recommendations
	if config.LearningRate < 0.1 {
		recommendations = append(recommendations, "Consider increasing learning rate for faster adaptation")
	} else if config.LearningRate > 0.8 {
		recommendations = append(recommendations, "Consider reducing learning rate for more stable learning")
	}

	// Boosting recommendations
	if config.BoostStrength == 0 && metrics.BoostingEvents == 0 {
		recommendations = append(recommendations, "Enable boosting to maintain column usage distribution")
	}

	// Performance recommendations
	if metrics.AverageProcessingTime > 5 {
		recommendations = append(recommendations, "Optimize configuration for faster processing")
	}

	return recommendations
}

// Private helper methods

func (h *SpatialPoolerHandler) validateProcessRequest(request *SpatialPoolerProcessRequest) error {
	// Validate required fields
	if request.InputID == "" {
		return fmt.Errorf("input_id is required")
	}

	if request.InputWidth <= 0 {
		return fmt.Errorf("input_width must be positive")
	}

	if request.EncoderOutput.Width <= 0 {
		return fmt.Errorf("encoder_output.width must be positive")
	}

	if len(request.EncoderOutput.ActiveBits) == 0 {
		return fmt.Errorf("encoder_output.active_bits cannot be empty")
	}

	if request.InputWidth != request.EncoderOutput.Width {
		return fmt.Errorf("input_width must match encoder_output.width")
	}

	// Validate active bits are within range
	for _, bit := range request.EncoderOutput.ActiveBits {
		if bit < 0 || bit >= request.EncoderOutput.Width {
			return fmt.Errorf("active bit %d out of range [0, %d)", bit, request.EncoderOutput.Width)
		}
	}

	return nil
}

// Request/Response types

type SpatialPoolerProcessRequest struct {
	EncoderOutput   EncoderOutputRequest   `json:"encoder_output" binding:"required"`
	InputWidth      int                    `json:"input_width" binding:"required,gt=0"`
	InputID         string                 `json:"input_id" binding:"required"`
	LearningEnabled bool                   `json:"learning_enabled"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

type EncoderOutputRequest struct {
	Width      int     `json:"width" binding:"required,gt=0"`
	ActiveBits []int   `json:"active_bits" binding:"required"`
	Sparsity   float64 `json:"sparsity" binding:"gte=0,lte=1"`
}

type SpatialPoolerProcessResponse struct {
	NormalizedSDR    SDRResponse `json:"normalized_sdr"`
	InputID          string      `json:"input_id"`
	ProcessingTimeMs int64       `json:"processing_time_ms"`
	ActiveColumns    []int       `json:"active_columns"`
	AvgOverlap       float64     `json:"avg_overlap"`
	SparsityLevel    float64     `json:"sparsity_level"`
	LearningOccurred bool        `json:"learning_occurred"`
	BoostingApplied  bool        `json:"boosting_applied"`
}

type SDRResponse struct {
	Width      int     `json:"width"`
	ActiveBits []int   `json:"active_bits"`
	Sparsity   float64 `json:"sparsity"`
}

type SpatialPoolerConfigResponse struct {
	InputWidth          int                        `json:"input_width"`
	ColumnCount         int                        `json:"column_count"`
	SparsityRatio       float64                    `json:"sparsity_ratio"`
	Mode                string                     `json:"mode"`
	LearningEnabled     bool                       `json:"learning_enabled"`
	LearningRate        float64                    `json:"learning_rate"`
	MaxBoost            float64                    `json:"max_boost"`
	BoostStrength       float64                    `json:"boost_strength"`
	InhibitionRadius    int                        `json:"inhibition_radius"`
	LocalAreaDensity    float64                    `json:"local_area_density"`
	MinOverlapThreshold int                        `json:"min_overlap_threshold"`
	MaxProcessingTimeMs int                        `json:"max_processing_time_ms"`
	SemanticThresholds  SemanticThresholdsResponse `json:"semantic_thresholds"`
}

type SemanticThresholdsResponse struct {
	SimilarInputMinOverlap   float64 `json:"similar_input_min_overlap"`
	DifferentInputMaxOverlap float64 `json:"different_input_max_overlap"`
}

type SpatialPoolerConfigUpdateRequest struct {
	InputWidth          int                             `json:"input_width" binding:"required,gt=0"`
	ColumnCount         int                             `json:"column_count" binding:"required,gt=0"`
	SparsityRatio       float64                         `json:"sparsity_ratio" binding:"required,gte=0.02,lte=0.05"`
	Mode                string                          `json:"mode" binding:"required"`
	LearningEnabled     bool                            `json:"learning_enabled"`
	LearningRate        float64                         `json:"learning_rate" binding:"gte=0,lte=1"`
	MaxBoost            float64                         `json:"max_boost" binding:"gte=1,lte=10"`
	BoostStrength       float64                         `json:"boost_strength" binding:"gte=0,lte=1"`
	InhibitionRadius    int                             `json:"inhibition_radius" binding:"gte=0"`
	LocalAreaDensity    float64                         `json:"local_area_density" binding:"gte=0,lte=1"`
	MinOverlapThreshold int                             `json:"min_overlap_threshold" binding:"gte=0"`
	MaxProcessingTimeMs int                             `json:"max_processing_time_ms" binding:"gte=1,lte=10"`
	SemanticThresholds  SemanticThresholdsUpdateRequest `json:"semantic_thresholds" binding:"required"`
}

type SemanticThresholdsUpdateRequest struct {
	SimilarInputMinOverlap   float64 `json:"similar_input_min_overlap" binding:"gte=0.3,lte=0.7"`
	DifferentInputMaxOverlap float64 `json:"different_input_max_overlap" binding:"gte=0,lte=0.2"`
}

type SpatialPoolerMetricsResponse struct {
	TotalProcessed           int64            `json:"total_processed"`
	AverageProcessingTimeMs  int64            `json:"average_processing_time_ms"`
	LearningIterations       int64            `json:"learning_iterations"`
	ColumnUsageDistribution  []float64        `json:"column_usage_distribution,omitempty"`
	AverageSparsity          float64          `json:"average_sparsity"`
	OverlapScoreDistribution []float64        `json:"overlap_score_distribution,omitempty"`
	BoostingEvents           int64            `json:"boosting_events"`
	ErrorCounts              map[string]int64 `json:"error_counts"`
}
