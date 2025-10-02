package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/htm-project/neural-api/internal/cortical/spatial"
	"github.com/htm-project/neural-api/internal/domain/htm"
	"github.com/htm-project/neural-api/internal/ports"
)

// spatialPoolingService implements the SpatialPoolingService interface
type spatialPoolingService struct {
	mu               sync.RWMutex
	engine           *spatial.SpatialPooler
	config           *htm.SpatialPoolerConfig
	observers        []ports.SpatialPoolingObserver
	instanceID       string
	createdAt        time.Time
	lastProcessingAt time.Time
}

// NewSpatialPoolingService creates a new spatial pooling service
func NewSpatialPoolingService(config *htm.SpatialPoolerConfig, instanceID string) (ports.SpatialPoolingService, error) {
	if config == nil {
		config = htm.DefaultSpatialPoolerConfig()
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create spatial pooler engine
	engine, err := spatial.NewSpatialPooler(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create spatial pooler: %w", err)
	}

	return &spatialPoolingService{
		engine:     engine,
		config:     config,
		observers:  make([]ports.SpatialPoolingObserver, 0),
		instanceID: instanceID,
		createdAt:  time.Now(),
	}, nil
}

// ProcessSpatialPooling transforms encoder output into normalized SDR
func (s *spatialPoolingService) ProcessSpatialPooling(ctx context.Context, input *htm.PoolingInput) (*htm.PoolingResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate input
	if err := input.Validate(); err != nil {
		s.notifyProcessingFailed(input.InputID, err)
		return nil, fmt.Errorf("input validation failed: %w", err)
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		err := ctx.Err()
		s.notifyProcessingFailed(input.InputID, err)
		return nil, err
	default:
	}

	// Notify processing started
	s.notifyProcessingStarted(input.InputID, input)

	// Process with spatial pooler
	result, err := s.engine.Process(input)
	if err != nil {
		s.notifyProcessingFailed(input.InputID, err)
		return nil, fmt.Errorf("spatial pooling failed: %w", err)
	}

	// Update last processing time
	s.lastProcessingAt = time.Now()

	// Notify processing completed
	s.notifyProcessingCompleted(input.InputID, result)

	return result, nil
}

// GetConfiguration returns current spatial pooler configuration
func (s *spatialPoolingService) GetConfiguration(ctx context.Context) (*htm.SpatialPoolerConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external modification
	configCopy := *s.config
	return &configCopy, nil
}

// UpdateConfiguration updates spatial pooler configuration
func (s *spatialPoolingService) UpdateConfiguration(ctx context.Context, config *htm.SpatialPoolerConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate new configuration
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Check for structural changes that require engine recreation
	if s.requiresEngineRecreation(config) {
		// Create new engine
		newEngine, err := spatial.NewSpatialPooler(config)
		if err != nil {
			return fmt.Errorf("failed to create new spatial pooler: %w", err)
		}

		oldConfig := s.config
		s.engine = newEngine
		s.config = config

		// Notify configuration changed
		s.notifyConfigurationChanged(oldConfig, config)
	} else {
		// Update existing engine
		if err := s.engine.UpdateConfiguration(config); err != nil {
			return fmt.Errorf("failed to update engine configuration: %w", err)
		}

		oldConfig := s.config
		s.config = config

		// Notify configuration changed
		s.notifyConfigurationChanged(oldConfig, config)
	}

	return nil
}

// GetMetrics returns spatial pooler performance metrics
func (s *spatialPoolingService) GetMetrics(ctx context.Context) (*htm.SpatialPoolerMetrics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.engine.GetMetrics(), nil
}

// ResetMetrics resets all performance metrics
func (s *spatialPoolingService) ResetMetrics(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Note: In a full implementation, we would reset the engine's metrics
	// For now, we'll document this as a limitation that would require
	// engine recreation or additional API in the spatial pooler

	return nil
}

// ValidateConfiguration validates a spatial pooler configuration
func (s *spatialPoolingService) ValidateConfiguration(ctx context.Context, config *htm.SpatialPoolerConfig) error {
	return config.Validate()
}

// HealthCheck performs health check on spatial pooler
func (s *spatialPoolingService) HealthCheck(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if engine exists
	if s.engine == nil {
		return fmt.Errorf("spatial pooler engine is not initialized")
	}

	// Check configuration validity
	if err := s.config.Validate(); err != nil {
		return fmt.Errorf("configuration is invalid: %w", err)
	}

	// Check processing pipeline with test input
	testInput := &htm.PoolingInput{
		EncoderOutput: htm.EncoderOutput{
			Width:      s.config.InputWidth,
			ActiveBits: []int{0, 1, 2}, // Simple test pattern
			Sparsity:   3.0 / float64(s.config.InputWidth),
		},
		InputWidth:      s.config.InputWidth,
		InputID:         "health-check",
		LearningEnabled: false,
	}

	// Test processing (without learning)
	_, err := s.engine.Process(testInput)
	if err != nil {
		return fmt.Errorf("health check processing failed: %w", err)
	}

	return nil
}

// GetInstanceInfo returns spatial pooler instance information
func (s *spatialPoolingService) GetInstanceInfo(ctx context.Context) map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info := map[string]interface{}{
		"instance_id":        s.instanceID,
		"created_at":         s.createdAt,
		"last_processing_at": s.lastProcessingAt,
		"uptime_seconds":     time.Since(s.createdAt).Seconds(),
		"configuration": map[string]interface{}{
			"input_width":      s.config.InputWidth,
			"column_count":     s.config.ColumnCount,
			"sparsity_ratio":   s.config.SparsityRatio,
			"mode":             s.config.Mode,
			"learning_enabled": s.config.LearningEnabled,
		},
		"observer_count": len(s.observers),
	}

	// Add metrics summary
	if metrics := s.engine.GetMetrics(); metrics != nil {
		info["metrics_summary"] = map[string]interface{}{
			"total_processed":         metrics.TotalProcessed,
			"average_processing_time": metrics.AverageProcessingTime,
			"average_sparsity":        metrics.AverageSparsity,
			"learning_iterations":     metrics.LearningIterations,
		}
	}

	return info
}

// AddObserver adds a processing observer
func (s *spatialPoolingService) AddObserver(observer ports.SpatialPoolingObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.observers = append(s.observers, observer)
}

// RemoveObserver removes a processing observer
func (s *spatialPoolingService) RemoveObserver(observer ports.SpatialPoolingObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Private helper methods

// requiresEngineRecreation checks if configuration changes require creating a new engine
func (s *spatialPoolingService) requiresEngineRecreation(newConfig *htm.SpatialPoolerConfig) bool {
	// Structural changes that require recreation
	if s.config.InputWidth != newConfig.InputWidth {
		return true
	}
	if s.config.ColumnCount != newConfig.ColumnCount {
		return true
	}

	return false
}

// Observer notification methods

func (s *spatialPoolingService) notifyProcessingStarted(inputID string, input *htm.PoolingInput) {
	for _, observer := range s.observers {
		observer.OnProcessingStarted(inputID, input)
	}
}

func (s *spatialPoolingService) notifyProcessingCompleted(inputID string, result *htm.PoolingResult) {
	for _, observer := range s.observers {
		observer.OnProcessingCompleted(inputID, result)
	}
}

func (s *spatialPoolingService) notifyProcessingFailed(inputID string, err error) {
	for _, observer := range s.observers {
		observer.OnProcessingFailed(inputID, err)
	}
}

func (s *spatialPoolingService) notifyConfigurationChanged(oldConfig, newConfig *htm.SpatialPoolerConfig) {
	for _, observer := range s.observers {
		observer.OnConfigurationChanged(oldConfig, newConfig)
	}
}

func (s *spatialPoolingService) notifyMetricsUpdated(metrics *htm.SpatialPoolerMetrics) {
	for _, observer := range s.observers {
		observer.OnMetricsUpdated(metrics)
	}
}

// SpatialPoolingServiceFactory creates spatial pooling services
type SpatialPoolingServiceFactory struct{}

// NewSpatialPoolingServiceFactory creates a new service factory
func NewSpatialPoolingServiceFactory() *SpatialPoolingServiceFactory {
	return &SpatialPoolingServiceFactory{}
}

// CreateService creates a spatial pooling service
func (f *SpatialPoolingServiceFactory) CreateService(config *htm.SpatialPoolerConfig, instanceID string) (ports.SpatialPoolingService, error) {
	return NewSpatialPoolingService(config, instanceID)
}

// CreateDefaultService creates a spatial pooling service with default configuration
func (f *SpatialPoolingServiceFactory) CreateDefaultService(instanceID string) (ports.SpatialPoolingService, error) {
	return NewSpatialPoolingService(htm.DefaultSpatialPoolerConfig(), instanceID)
}

// ValidateServiceConfiguration validates service configuration
func (f *SpatialPoolingServiceFactory) ValidateServiceConfiguration(config *htm.SpatialPoolerConfig) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}
	return config.Validate()
}
