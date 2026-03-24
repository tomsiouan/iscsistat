package config

import (
	"fmt"
	"strings"
	"time"

	"iscsistat/internal/validator"
)

// Config represents the application's global configuration.
type Config struct {
	HTTP    HTTPConfig
	Metrics MetricsConfig
}

// HTTPConfig holds the configuration settings for the HTTP server.
type HTTPConfig struct {
	Host string
	Port int
	TLS  TLSConfig
}

// TLSConfig defines the parameters for Transport Layer Security.
type TLSConfig struct {
	Enabled      bool
	CertFile     string
	KeyFile      string
	ClientCAPath string
}

type MetricsConfig struct {
	CollectInterval time.Duration
}

// LoadConfig loads configuration from environment variables and returns a Config.
// It returns an error if loading or validation fails.
func LoadConfig() (Config, error) {
	cfg := Config{
		HTTP: HTTPConfig{
			Host: env("HTTP_HOST", "0.0.0.0"),
			Port: envInt("HTTP_PORT", 9101),
			TLS: TLSConfig{
				Enabled:      envBool("HTTP_TLS_ENABLED", false),
				CertFile:     env("HTTP_TLS_CERT_FILE", ""),
				KeyFile:      env("HTTP_TLS_KEY_FILE", ""),
				ClientCAPath: env("HTTP_TLS_CLIENT_CA_FILE", ""),
			},
		},
		Metrics: MetricsConfig{
			CollectInterval: envDuration("METRICS_COLLECT_INTERVAL", 15*time.Second),
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate checks the configuration for consistency and required fields.
// It returns a formatted error if any validation rules are violated, nil otherwise.
func (c Config) Validate() error {
	v := validator.New()

	if c.HTTP.Port < 1 || c.HTTP.Port > 65535 {
		v.Add("HTTP_PORT", fmt.Sprintf("%d is out of range", c.HTTP.Port))
	}

	if c.HTTP.TLS.Enabled {
		if c.HTTP.TLS.CertFile == "" {
			v.Add("HTTP_TLS_CERT_FILE", "required when TLS is enabled")
		}
		if c.HTTP.TLS.KeyFile == "" {
			v.Add("HTTP_TLS_KEY_FILE", "required when TLS is enabled")
		}
	}

	if c.Metrics.CollectInterval <= 0 {
		v.Add("METRICS_COLLECT_INTERVAL", "must be a positive duration (e.g., 15s, 1m)")
	}

	if !v.IsValid() {
		return formatValidationErrors(v.Errors)
	}

	return nil
}

// formatValidationErrors aggregates a slice of FieldError into a single human-readable error.
func formatValidationErrors(errors []validator.FieldError) error {
	var sb strings.Builder
	sb.WriteString("configuration errors:")
	for _, err := range errors {
		sb.WriteString(fmt.Sprintf("\n - %s: %s", err.Field, err.Message))
	}
	return fmt.Errorf(sb.String())
}
