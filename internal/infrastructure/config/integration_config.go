package config

import (
	"fmt"
	"time"

	"github.com/htm-project/neural-api/internal/domain/htm"
)

// IntegrationConfig manages the complete integration configuration
type IntegrationConfig struct {
	Application *ApplicationConfig        `json:"application" validate:"required"`
	Server      *IntegrationServerConfig  `json:"server" validate:"required"`
	Performance *PerformanceConfig        `json:"performance" validate:"required"`
	Logging     *IntegrationLoggingConfig `json:"logging" validate:"required"`
	Metrics     *IntegrationMetricsConfig `json:"metrics" validate:"required"`
}

// ApplicationConfig holds complete application configuration
type ApplicationConfig struct {
	Name          string                   `json:"name" validate:"required"`
	Version       string                   `json:"version" validate:"required"`
	Environment   string                   `json:"environment" validate:"required,oneof=development production testing"`
	SpatialPooler *htm.SpatialPoolerConfig `json:"spatial_pooler" validate:"required"`
	Features      *FeatureConfig           `json:"features" validate:"required"`
}

// IntegrationServerConfig holds HTTP server configuration for production deployment
type IntegrationServerConfig struct {
	Host            string        `json:"host" validate:"required"`
	Port            int           `json:"port" validate:"required,min=1,max=65535"`
	ReadTimeout     time.Duration `json:"read_timeout" validate:"required"`
	WriteTimeout    time.Duration `json:"write_timeout" validate:"required"`
	IdleTimeout     time.Duration `json:"idle_timeout" validate:"required"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout" validate:"required"`
	MaxHeaderBytes  int           `json:"max_header_bytes" validate:"required,min=1"`
	EnableProfiling bool          `json:"enable_profiling"`
	EnableCORS      bool          `json:"enable_cors"`
	TrustedProxies  []string      `json:"trusted_proxies"`
}

// Address returns the server address in host:port format
func (s *IntegrationServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// PerformanceConfig holds performance tuning parameters
type PerformanceConfig struct {
	MaxConcurrentRequests int  `json:"max_concurrent_requests" validate:"required,min=1,max=1000"`
	MaxDatasetSizeMB      int  `json:"max_dataset_size_mb" validate:"required,min=1,max=100"`
	ResponseTimeTargetMs  int  `json:"response_time_target_ms" validate:"required,min=1,max=1000"`
	MatrixPoolSize        int  `json:"matrix_pool_size" validate:"required,min=1,max=100"`
	GCTargetPercentage    int  `json:"gc_target_percentage" validate:"required,min=10,max=500"`
	MemoryLimitMB         int  `json:"memory_limit_mb" validate:"required,min=100,max=10000"`
	RequestTimeoutMs      int  `json:"request_timeout_ms" validate:"required,min=100,max=30000"`
	EnableMatrixPooling   bool `json:"enable_matrix_pooling"`
	EnableGCOptimization  bool `json:"enable_gc_optimization"`
	EnableRequestLimiting bool `json:"enable_request_limiting"`
}

// IntegrationLoggingConfig holds logging configuration
type IntegrationLoggingConfig struct {
	Level            string `json:"level" validate:"required,oneof=debug info warn error"`
	Format           string `json:"format" validate:"required,oneof=json text"`
	EnableStructured bool   `json:"enable_structured"`
	EnableHTMContext bool   `json:"enable_htm_context"`
	LogRequests      bool   `json:"log_requests"`
	LogResponses     bool   `json:"log_responses"`
	LogPerformance   bool   `json:"log_performance"`
}

// IntegrationMetricsConfig holds metrics collection configuration
type IntegrationMetricsConfig struct {
	Enabled            bool          `json:"enabled"`
	CollectionInterval time.Duration `json:"collection_interval" validate:"required"`
	RetentionPeriod    time.Duration `json:"retention_period" validate:"required"`
	EnableHTMMetrics   bool          `json:"enable_htm_metrics"`
	EnablePerformance  bool          `json:"enable_performance"`
	EnableHealthChecks bool          `json:"enable_health_checks"`
	ExportPrometheus   bool          `json:"export_prometheus"`
	PrometheusPort     int           `json:"prometheus_port" validate:"min=1,max=65535"`
}

// FeatureConfig holds feature flags and configuration
type FeatureConfig struct {
	EnableDebugEndpoints     bool `json:"enable_debug_endpoints"`
	EnableMetricsEndpoints   bool `json:"enable_metrics_endpoints"`
	EnableValidationCache    bool `json:"enable_validation_cache"`
	EnableSparsityValidation bool `json:"enable_sparsity_validation"`
	EnableHTMValidation      bool `json:"enable_htm_validation"`
}

// NewDefaultIntegrationConfig creates a default integration configuration
func NewDefaultIntegrationConfig() *IntegrationConfig {
	return &IntegrationConfig{
		Application: &ApplicationConfig{
			Name:        "HTM Spatial Pooler API",
			Version:     "1.0.0",
			Environment: "production",
			SpatialPooler: &htm.SpatialPoolerConfig{
				InputWidth:          1024,
				ColumnCount:         2048,
				SparsityRatio:       0.025, // Increased to 2.5% to ensure we stay above 2% minimum
				Mode:                htm.SpatialPoolerModeDeterministic,
				LearningEnabled:     true,
				LearningRate:        0.1,
				MaxBoost:            3.0,
				BoostStrength:       0.5,
				InhibitionRadius:    16,
				LocalAreaDensity:    0.025, // Match sparsity ratio
				MinOverlapThreshold: 5,
				MaxProcessingTimeMs: 50, // Increased to 50ms for realistic performance expectation
				SemanticThresholds: htm.SemanticThresholds{
					SimilarInputMinOverlap:   0.5,
					DifferentInputMaxOverlap: 0.1,
				},
			},
			Features: &FeatureConfig{
				EnableDebugEndpoints:     false,
				EnableMetricsEndpoints:   true,
				EnableValidationCache:    true,
				EnableSparsityValidation: true,
				EnableHTMValidation:      true,
			},
		},
		Server: &IntegrationServerConfig{
			Host:            "0.0.0.0",
			Port:            8080,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			IdleTimeout:     120 * time.Second,
			ShutdownTimeout: 30 * time.Second,
			MaxHeaderBytes:  1 << 20, // 1MB
			EnableProfiling: false,
			EnableCORS:      true,
			TrustedProxies:  []string{"127.0.0.1", "::1"},
		},
		Performance: &PerformanceConfig{
			MaxConcurrentRequests: 100,
			MaxDatasetSizeMB:      10,
			ResponseTimeTargetMs:  100,
			MatrixPoolSize:        10,
			GCTargetPercentage:    100,
			MemoryLimitMB:         512,
			RequestTimeoutMs:      5000,
			EnableMatrixPooling:   true,
			EnableGCOptimization:  true,
			EnableRequestLimiting: true,
		},
		Logging: &IntegrationLoggingConfig{
			Level:            "info",
			Format:           "json",
			EnableStructured: true,
			EnableHTMContext: true,
			LogRequests:      true,
			LogResponses:     false,
			LogPerformance:   true,
		},
		Metrics: &IntegrationMetricsConfig{
			Enabled:            true,
			CollectionInterval: 30 * time.Second,
			RetentionPeriod:    24 * time.Hour,
			EnableHTMMetrics:   true,
			EnablePerformance:  true,
			EnableHealthChecks: true,
			ExportPrometheus:   true,
			PrometheusPort:     9090,
		},
	}
}

// Validate performs comprehensive validation of the integration configuration
func (ic *IntegrationConfig) Validate() error {
	if ic.Application == nil {
		return fmt.Errorf("application configuration is required")
	}

	if ic.Server == nil {
		return fmt.Errorf("server configuration is required")
	}

	if ic.Performance == nil {
		return fmt.Errorf("performance configuration is required")
	}

	if ic.Logging == nil {
		return fmt.Errorf("logging configuration is required")
	}

	if ic.Metrics == nil {
		return fmt.Errorf("metrics configuration is required")
	}

	// Validate application configuration
	if err := ic.validateApplication(); err != nil {
		return fmt.Errorf("application config validation failed: %w", err)
	}

	// Validate server configuration
	if err := ic.validateServer(); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	// Validate performance configuration
	if err := ic.validatePerformance(); err != nil {
		return fmt.Errorf("performance config validation failed: %w", err)
	}

	return nil
}

// validateApplication validates application configuration
func (ic *IntegrationConfig) validateApplication() error {
	app := ic.Application

	if app.Name == "" {
		return fmt.Errorf("application name is required")
	}

	if app.Version == "" {
		return fmt.Errorf("application version is required")
	}

	if app.SpatialPooler == nil {
		return fmt.Errorf("spatial pooler configuration is required")
	}

	// Validate spatial pooler configuration
	if app.SpatialPooler.SparsityRatio <= 0 || app.SpatialPooler.SparsityRatio > 0.1 {
		return fmt.Errorf("spatial pooler sparsity ratio must be between 0 and 0.1")
	}

	if app.SpatialPooler.ColumnCount <= 0 {
		return fmt.Errorf("spatial pooler column count must be positive")
	}

	if app.SpatialPooler.InputWidth <= 0 {
		return fmt.Errorf("spatial pooler input width must be positive")
	}

	return nil
}

// validateServer validates server configuration
func (ic *IntegrationConfig) validateServer() error {
	server := ic.Server

	if server.Port <= 0 || server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}

	if server.ReadTimeout <= 0 {
		return fmt.Errorf("server read timeout must be positive")
	}

	if server.WriteTimeout <= 0 {
		return fmt.Errorf("server write timeout must be positive")
	}

	return nil
}

// validatePerformance validates performance configuration
func (ic *IntegrationConfig) validatePerformance() error {
	perf := ic.Performance

	if perf.MaxConcurrentRequests <= 0 || perf.MaxConcurrentRequests > 1000 {
		return fmt.Errorf("max concurrent requests must be between 1 and 1000")
	}

	if perf.MaxDatasetSizeMB <= 0 || perf.MaxDatasetSizeMB > 100 {
		return fmt.Errorf("max dataset size must be between 1 and 100 MB")
	}

	if perf.ResponseTimeTargetMs <= 0 || perf.ResponseTimeTargetMs > 1000 {
		return fmt.Errorf("response time target must be between 1 and 1000 ms")
	}

	return nil
}
