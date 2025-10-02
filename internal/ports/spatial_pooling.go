package ports

import (
	"context"

	"github.com/htm-project/neural-api/internal/domain/htm"
)

// SpatialPoolingService defines the interface for spatial pooling operations
type SpatialPoolingService interface {
	// ProcessSpatialPooling transforms encoder output into normalized SDR
	ProcessSpatialPooling(ctx context.Context, input *htm.PoolingInput) (*htm.PoolingResult, error)

	// GetConfiguration returns current spatial pooler configuration
	GetConfiguration(ctx context.Context) (*htm.SpatialPoolerConfig, error)

	// UpdateConfiguration updates spatial pooler configuration
	UpdateConfiguration(ctx context.Context, config *htm.SpatialPoolerConfig) error

	// GetMetrics returns spatial pooler performance metrics
	GetMetrics(ctx context.Context) (*htm.SpatialPoolerMetrics, error)

	// ResetMetrics resets all performance metrics
	ResetMetrics(ctx context.Context) error

	// ValidateConfiguration validates a spatial pooler configuration
	ValidateConfiguration(ctx context.Context, config *htm.SpatialPoolerConfig) error

	// HealthCheck performs health check on spatial pooler
	HealthCheck(ctx context.Context) error

	// GetInstanceInfo returns spatial pooler instance information
	GetInstanceInfo(ctx context.Context) map[string]interface{}
}

// SpatialPoolingPort defines the port interface for spatial pooling infrastructure
type SpatialPoolingPort interface {
	// CreateSpatialPooler creates a new spatial pooler instance
	CreateSpatialPooler(config *htm.SpatialPoolerConfig) (SpatialPoolingEngine, error)

	// GetDefaultConfiguration returns default spatial pooler configuration
	GetDefaultConfiguration() *htm.SpatialPoolerConfig

	// ValidateInput validates spatial pooling input
	ValidateInput(input *htm.PoolingInput) error

	// ValidateOutput validates spatial pooling output
	ValidateOutput(result *htm.PoolingResult) error
}

// SpatialPoolingEngine defines the core spatial pooling computation engine
type SpatialPoolingEngine interface {
	// Process performs spatial pooling transformation
	Process(input *htm.PoolingInput) (*htm.PoolingResult, error)

	// GetConfiguration returns current configuration
	GetConfiguration() *htm.SpatialPoolerConfig

	// UpdateConfiguration updates configuration (only non-structural changes)
	UpdateConfiguration(config *htm.SpatialPoolerConfig) error

	// GetMetrics returns current metrics
	GetMetrics() *htm.SpatialPoolerMetrics

	// ResetMetrics resets performance metrics
	ResetMetrics()

	// IsHealthy returns true if engine is operating normally
	IsHealthy() bool

	// GetDiagnostics returns diagnostic information
	GetDiagnostics() map[string]interface{}
}

// SpatialPoolingRepository defines the interface for persistence (if needed)
type SpatialPoolingRepository interface {
	// SaveConfiguration persists spatial pooler configuration
	SaveConfiguration(instanceID string, config *htm.SpatialPoolerConfig) error

	// LoadConfiguration loads spatial pooler configuration
	LoadConfiguration(instanceID string) (*htm.SpatialPoolerConfig, error)

	// SaveMetrics persists spatial pooler metrics
	SaveMetrics(instanceID string, metrics *htm.SpatialPoolerMetrics) error

	// LoadMetrics loads spatial pooler metrics
	LoadMetrics(instanceID string) (*htm.SpatialPoolerMetrics, error)

	// DeleteInstance removes all data for a spatial pooler instance
	DeleteInstance(instanceID string) error
}

// SpatialPoolingObserver defines the interface for monitoring spatial pooling operations
type SpatialPoolingObserver interface {
	// OnProcessingStarted is called when spatial pooling processing begins
	OnProcessingStarted(inputID string, input *htm.PoolingInput)

	// OnProcessingCompleted is called when spatial pooling processing completes
	OnProcessingCompleted(inputID string, result *htm.PoolingResult)

	// OnProcessingFailed is called when spatial pooling processing fails
	OnProcessingFailed(inputID string, err error)

	// OnConfigurationChanged is called when configuration is updated
	OnConfigurationChanged(oldConfig, newConfig *htm.SpatialPoolerConfig)

	// OnMetricsUpdated is called when metrics are updated
	OnMetricsUpdated(metrics *htm.SpatialPoolerMetrics)
}

// SpatialPoolingFactory defines the interface for creating spatial pooling components
type SpatialPoolingFactory interface {
	// CreateService creates a spatial pooling service
	CreateService(config *htm.SpatialPoolerConfig) (SpatialPoolingService, error)

	// CreateEngine creates a spatial pooling engine
	CreateEngine(config *htm.SpatialPoolerConfig) (SpatialPoolingEngine, error)

	// CreateRepository creates a spatial pooling repository
	CreateRepository() SpatialPoolingRepository

	// CreateObserver creates a spatial pooling observer
	CreateObserver() SpatialPoolingObserver
}

// SpatialPoolingAdapter defines the interface for adapting external spatial pooling implementations
type SpatialPoolingAdapter interface {
	// Adapt wraps an external spatial pooling implementation
	Adapt(externalEngine interface{}) (SpatialPoolingEngine, error)

	// IsCompatible checks if external implementation is compatible
	IsCompatible(externalEngine interface{}) bool

	// GetRequiredInterface returns the required interface for external engines
	GetRequiredInterface() interface{}
}
