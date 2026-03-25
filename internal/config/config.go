package config

import (
	"fmt"
	"strings"
	"time"

	"iscsistat/internal/iscsi"
	"iscsistat/internal/validator"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

// Config represents the application's global configuration.
type Config struct {
	HTTP    HTTPConfig    `toml:"http"`
	Metrics MetricsConfig `toml:"metrics"`
	ISCSI   ISCSIConfig   `toml:"iscsi"`
}

// HTTPConfig holds the configuration settings for the HTTP server.
type HTTPConfig struct {
	Host string    `toml:"host"`
	Port int       `toml:"port"`
	TLS  TLSConfig `toml:"tls"`
}

// TLSConfig defines the parameters for Transport Layer Security.
type TLSConfig struct {
	Enabled      bool   `toml:"enabled"`
	CertFile     string `toml:"cert_file"`
	KeyFile      string `toml:"key_file"`
	ClientCAPath string `toml:"client_ca_path"`
}

// MetricsConfig defines the settings related to the metrics collection process.
type MetricsConfig struct {
	CollectInterval time.Duration `toml:"collect_interval"`
}

type ISCSIConfig struct {
	Parser     string `toml:"parser"`
	TrimPrefix string `toml:"trim_prefix"`
	TrimSuffix string `toml:"trim_suffix"`
}

// LoadConfig acts as the single entry point for application configuration.
// It follows a priority logic: if a file path is provided, it loads from TOML;
// otherwise, it falls back to loading values from environment variables.
func LoadConfig(path string, logger *zap.Logger) (Config, error) {
	if path != "" {
		return loadFromFile(path, logger)
	}

	return loadFromEnv(logger)
}

// loadFromFile reads, parses, and validates the configuration from a TOML file.
// It returns a Config struct or an error if the file is missing, malformed,
// or contains invalid values.
func loadFromFile(filePath string, logger *zap.Logger) (Config, error) {
	var cfg Config

	if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse TOML config file: %w", err)
	}

	if cfg.ISCSI.Parser == "" {
		cfg.ISCSI.Parser = iscsi.DefaultParser
	}

	if err := cfg.Validate(logger); err != nil {
		return Config{}, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// LoadConfig loads configuration from environment variables and returns a Config.
// It returns an error if loading or validation fails.
func loadFromEnv(logger *zap.Logger) (Config, error) {
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
		ISCSI: ISCSIConfig{
			Parser: env("ISCSI_PARSER", iscsi.DefaultParser),
		},
	}

	if err := cfg.Validate(logger); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate checks the configuration for consistency and required fields.
// It returns a formatted error if any validation rules are violated, nil otherwise.
func (c Config) Validate(logger *zap.Logger) error {
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

	if !iscsi.ParserExists(c.ISCSI.Parser) {
		v.Add("ISCSI_PARSER", fmt.Sprintf("'%s' is not a valid parser", c.ISCSI.Parser))
	}

	if !v.IsValid() {
		return formatValidationErrors(v.Errors)
	}

	logger.Info("Configuration loaded successfully",
		zap.String("parser", c.ISCSI.Parser),
	)

	return nil
}

// formatValidationErrors aggregates a slice of FieldError into a single human-readable error.
func formatValidationErrors(errors []validator.FieldError) error {
	var sb strings.Builder
	sb.WriteString("configuration errors:")
	for _, err := range errors {
		sb.WriteString(fmt.Sprintf("\n - %s: %s", err.Field, err.Message))
	}
	return fmt.Errorf("%s", sb.String())
}
