package handlers

import (
	"context"
	"runtime"
	"time"

	"github.com/htm-project/neural-api/internal/ports"
)

// HealthHandlerImpl implements the HealthHandler interface.
type HealthHandlerImpl struct {
	processingService     ports.ProcessingService
	spatialPoolingService ports.SpatialPoolingService
	metricsCollector      ports.MetricsCollector
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(
	processingService ports.ProcessingService,
	spatialPoolingService ports.SpatialPoolingService,
	metricsCollector ports.MetricsCollector,
) ports.HealthHandler {
	return &HealthHandlerImpl{
		processingService:     processingService,
		spatialPoolingService: spatialPoolingService,
		metricsCollector:      metricsCollector,
	}
}

// HandleHealthCheck performs health checks and returns status.
func (h *HealthHandlerImpl) HandleHealthCheck(ctx context.Context) (map[string]interface{}, error) {
	healthData := make(map[string]interface{})

	// Get dependency health status
	dependencies := h.CheckDependencies(ctx)
	healthData["dependencies"] = dependencies

	// Get system information
	systemInfo := h.GetSystemInfo()
	healthData["system"] = systemInfo

	// Add service-specific health checks
	serviceHealth := map[string]interface{}{
		"processing_service":      h.processingService != nil,
		"spatial_pooling_service": h.spatialPoolingService != nil,
		"metrics_collector":       h.metricsCollector != nil,
		"uptime_seconds":          time.Since(startTime).Seconds(),
	}

	// Check spatial pooler health
	spatialPoolerHealthy := true
	if h.spatialPoolingService != nil {
		if err := h.spatialPoolingService.HealthCheck(ctx); err != nil {
			spatialPoolerHealthy = false
			serviceHealth["spatial_pooler_error"] = err.Error()
		}
	} else {
		spatialPoolerHealthy = false
	}
	serviceHealth["spatial_pooler_healthy"] = spatialPoolerHealthy

	healthData["service"] = serviceHealth

	// Overall health status
	allHealthy := true
	for _, healthy := range dependencies {
		if !healthy {
			allHealthy = false
			break
		}
	}

	// Include spatial pooler health in overall status
	allHealthy = allHealthy && spatialPoolerHealthy

	healthData["healthy"] = allHealthy

	return healthData, nil
}

// CheckDependencies checks the health of all dependencies.
func (h *HealthHandlerImpl) CheckDependencies(ctx context.Context) map[string]bool {
	dependencies := make(map[string]bool)

	// Check processing service
	dependencies["processing_service"] = h.checkProcessingService(ctx)

	// Check metrics collector
	dependencies["metrics_collector"] = h.checkMetricsCollector(ctx)

	// Check memory usage
	dependencies["memory"] = h.checkMemoryUsage()

	// Check disk space (simplified check)
	dependencies["disk"] = h.checkDiskSpace()

	return dependencies
}

// GetSystemInfo returns basic system information.
func (h *HealthHandlerImpl) GetSystemInfo() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]interface{}{
		"go_version":     runtime.Version(),
		"go_os":          runtime.GOOS,
		"go_arch":        runtime.GOARCH,
		"num_cpu":        runtime.NumCPU(),
		"num_goroutines": runtime.NumGoroutine(),
		"memory": map[string]interface{}{
			"alloc_mb":       bytesToMB(memStats.Alloc),
			"total_alloc_mb": bytesToMB(memStats.TotalAlloc),
			"sys_mb":         bytesToMB(memStats.Sys),
			"num_gc":         memStats.NumGC,
		},
		"build_info": map[string]interface{}{
			"version": "1.0.0",
			"commit":  "dev",
			"date":    "2024-01-01",
		},
	}
}

// checkProcessingService checks if the processing service is healthy.
func (h *HealthHandlerImpl) checkProcessingService(ctx context.Context) bool {
	if h.processingService == nil {
		return false
	}

	// Try to perform a simple health check operation
	// This would depend on the actual processing service implementation
	// For now, we just check if it's not nil
	return true
}

// checkMetricsCollector checks if the metrics collector is healthy.
func (h *HealthHandlerImpl) checkMetricsCollector(ctx context.Context) bool {
	if h.metricsCollector == nil {
		return false
	}

	// Try to record a test metric
	// This verifies the metrics collector is functioning
	h.metricsCollector.IncrementRequestCount()
	return true
}

// checkMemoryUsage checks if memory usage is within acceptable limits.
func (h *HealthHandlerImpl) checkMemoryUsage() bool {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Check if allocated memory is less than 1GB
	const maxMemoryMB = 1024
	allocMB := bytesToMB(memStats.Alloc)

	return allocMB < maxMemoryMB
}

// checkDiskSpace checks if there's sufficient disk space.
func (h *HealthHandlerImpl) checkDiskSpace() bool {
	// Simplified disk space check
	// In a real implementation, this would check actual disk usage
	return true
}

// MetricsHandlerImpl implements the MetricsHandler interface.
type MetricsHandlerImpl struct {
	metricsCollector ports.MetricsCollector
}

// NewMetricsHandler creates a new metrics handler.
func NewMetricsHandler(metricsCollector ports.MetricsCollector) ports.MetricsHandler {
	return &MetricsHandlerImpl{
		metricsCollector: metricsCollector,
	}
}

// HandleMetrics returns current system metrics.
func (m *MetricsHandlerImpl) HandleMetrics(ctx context.Context) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	// Get performance metrics
	performanceMetrics := m.GetPerformanceMetrics()
	metrics["performance"] = performanceMetrics

	// Get request metrics
	requestMetrics := m.GetRequestMetrics()
	metrics["requests"] = requestMetrics

	// Get system metrics
	systemMetrics := m.GetSystemMetrics()
	metrics["system"] = systemMetrics

	// Add collector-specific metrics if available
	if m.metricsCollector != nil {
		// This would depend on the actual metrics collector implementation
		// For now, we add a placeholder
		metrics["collector_status"] = "active"
	}

	return metrics, nil
}

// GetPerformanceMetrics returns performance-related metrics.
func (m *MetricsHandlerImpl) GetPerformanceMetrics() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]interface{}{
		"memory": map[string]interface{}{
			"heap_alloc_mb":    bytesToMB(memStats.HeapAlloc),
			"heap_sys_mb":      bytesToMB(memStats.HeapSys),
			"heap_idle_mb":     bytesToMB(memStats.HeapIdle),
			"heap_inuse_mb":    bytesToMB(memStats.HeapInuse),
			"heap_released_mb": bytesToMB(memStats.HeapReleased),
			"heap_objects":     memStats.HeapObjects,
		},
		"gc": map[string]interface{}{
			"num_gc":          memStats.NumGC,
			"pause_total_ns":  memStats.PauseTotalNs,
			"gc_cpu_fraction": memStats.GCCPUFraction,
		},
		"goroutines":     runtime.NumGoroutine(),
		"uptime_seconds": time.Since(startTime).Seconds(),
	}
}

// GetRequestMetrics returns request-related metrics.
func (m *MetricsHandlerImpl) GetRequestMetrics() map[string]interface{} {
	// These would come from the actual metrics collector
	// For now, we return placeholder data
	return map[string]interface{}{
		"total_requests":           0,
		"successful_requests":      0,
		"failed_requests":          0,
		"average_response_time_ms": 0,
		"requests_per_second":      0,
		"active_requests":          0,
	}
}

// GetSystemMetrics returns system-related metrics.
func (m *MetricsHandlerImpl) GetSystemMetrics() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]interface{}{
		"cpu_count":  runtime.NumCPU(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"go_version": runtime.Version(),
		"memory": map[string]interface{}{
			"sys_mb":         bytesToMB(memStats.Sys),
			"total_alloc_mb": bytesToMB(memStats.TotalAlloc),
			"mallocs":        memStats.Mallocs,
			"frees":          memStats.Frees,
		},
		"timestamps": map[string]interface{}{
			"start_time":   startTime.Format(time.RFC3339),
			"current_time": time.Now().Format(time.RFC3339),
		},
	}
}

// Utility functions

var startTime = time.Now()

// bytesToMB converts bytes to megabytes.
func bytesToMB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}
