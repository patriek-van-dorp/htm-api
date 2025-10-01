package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration.
type Config struct {
	Port        int    `json:"port"`
	Environment string `json:"environment"`
	Debug       bool   `json:"debug"`
	Timeout     int    `json:"timeout"`
}

// Load loads configuration from environment variables with defaults.
func Load() (*Config, error) {
	cfg := &Config{
		Port:        8080,
		Environment: "development",
		Debug:       true,
		Timeout:     30,
	}

	// Load from environment variables
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		cfg.Environment = env
	}

	if env := os.Getenv("ENV"); env != "" {
		cfg.Environment = env
	}

	if debug := os.Getenv("DEBUG"); debug != "" {
		if d, err := strconv.ParseBool(debug); err == nil {
			cfg.Debug = d
		}
	}

	if timeout := os.Getenv("TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.Timeout = t
		}
	}

	// Set debug based on environment
	if cfg.Environment == "production" {
		cfg.Debug = false
	}

	return cfg, nil
}
