package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the HTM Neural API
type Config struct {
	Server  ServerConfig
	API     APIConfig
	Logging LoggingConfig
	Metrics MetricsConfig
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// APIConfig contains API-specific configuration
type APIConfig struct {
	Version                         string
	MaxRequestSize                  int64
	DefaultProcessingTimeoutTimeout time.Duration
	MaxConcurrentRequests           int
	EnableCORS                      bool
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string
	Format string // "json" or "text"
}

// MetricsConfig contains metrics collection configuration
type MetricsConfig struct {
	Enabled bool
	Path    string
}

// Load reads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "localhost"),
			Port:            getEnv("SERVER_PORT", "8080"),
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		API: APIConfig{
			Version:                         getEnv("API_VERSION", "v1.0"),
			MaxRequestSize:                  getIntEnv("API_MAX_REQUEST_SIZE", 10*1024*1024), // 10MB
			DefaultProcessingTimeoutTimeout: getDurationEnv("API_PROCESSING_TIMEOUT", 5*time.Minute),
			MaxConcurrentRequests:           int(getIntEnv("API_MAX_CONCURRENT_REQUESTS", 100)),
			EnableCORS:                      getBoolEnv("API_ENABLE_CORS", true),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Metrics: MetricsConfig{
			Enabled: getBoolEnv("METRICS_ENABLED", true),
			Path:    getEnv("METRICS_PATH", "/metrics"),
		},
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv gets an integer environment variable with a default value
func getIntEnv(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getBoolEnv gets a boolean environment variable with a default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getDurationEnv gets a duration environment variable with a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// Address returns the full server address for binding
func (c *ServerConfig) Address() string {
	return c.Host + ":" + c.Port
}
