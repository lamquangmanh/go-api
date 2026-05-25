package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config is the root configuration struct, representing the full configs/config.yml file
type Config struct {
	GRPC     GRPCConfig     `yaml:"grpc"`
	HTTP     HTTPConfig     `yaml:"http"`
	Database DatabaseConfig `yaml:"database"`
	Logger   LoggerConfig   `yaml:"logger"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// GRPCConfig holds gRPC server configuration
type GRPCConfig struct {
	Port int `yaml:"port"`
}

// HTTPConfig holds HTTP server configuration (used for /healthz endpoint)
type HTTPConfig struct {
	Port int `yaml:"port"`
}

// DatabaseConfig holds PostgreSQL connection parameters
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

// LoggerConfig holds logging behavior configuration.
// Level supports: debug, info, warn, error.
// Format supports: json, console.
type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// JWTConfig holds JWT configuration.
type JWTConfig struct {
	Secret 				string `yaml:"secret"`
	ExpirySeconds int  `yaml:"expiry_seconds"`
}

// applyDefaults fills optional config values with sensible defaults.
// This keeps startup behavior predictable even when fields are omitted.
func (c *Config) applyDefaults() {
	if c.GRPC.Port == 0 {
		c.GRPC.Port = 5100
	}
	if c.HTTP.Port == 0 {
		c.HTTP.Port = 5101
	}

	if c.Database.Host == "" {
		c.Database.Host = "localhost"
	}
	if c.Database.Port == 0 {
		c.Database.Port = 5432
	}
	if c.Database.User == "" {
		c.Database.User = "postgres"
	}
	if c.Database.Name == "" {
		c.Database.Name = "go_api"
	}
	if c.Database.SSLMode == "" {
		c.Database.SSLMode = "disable"
	}

	if c.Logger.Level == "" {
		c.Logger.Level = "info"
	}
	if c.Logger.Format == "" {
		c.Logger.Format = "json"
	}

	if c.JWT.Secret == "" {
		c.JWT.Secret = "go-api-dev-secret"
	}
	if c.JWT.ExpirySeconds == 0 {
		c.JWT.ExpirySeconds = 86400 // 1 day
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func applyIntEnv(current int, keys ...string) int {
	raw := firstNonEmpty(getEnvValues(keys...)...)
	if raw == "" {
		return current
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return current
	}
	return parsed
}

func applyStringEnv(current string, keys ...string) string {
	raw := firstNonEmpty(getEnvValues(keys...)...)
	if raw == "" {
		return current
	}
	return raw
}

func getEnvValues(keys ...string) []string {
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, os.Getenv(key))
	}
	return values
}

func (c *Config) applyEnvOverrides() {
	c.GRPC.Port = applyIntEnv(c.GRPC.Port, "GOAPI_GRPC_PORT", "GRPC_PORT")
	c.HTTP.Port = applyIntEnv(c.HTTP.Port, "GOAPI_HTTP_PORT", "HTTP_PORT")

	c.Database.Host = applyStringEnv(c.Database.Host, "GOAPI_DB_HOST", "DB_HOST")
	c.Database.Port = applyIntEnv(c.Database.Port, "GOAPI_DB_PORT", "DB_PORT")
	c.Database.User = applyStringEnv(c.Database.User, "GOAPI_DB_USER", "DB_USER")
	c.Database.Password = applyStringEnv(c.Database.Password, "GOAPI_DB_PASSWORD", "DB_PASSWORD")
	c.Database.Name = applyStringEnv(c.Database.Name, "GOAPI_DB_NAME", "DB_NAME")
	c.Database.SSLMode = applyStringEnv(c.Database.SSLMode, "GOAPI_DB_SSL_MODE", "DB_SSL_MODE")

	c.Logger.Level = applyStringEnv(c.Logger.Level, "GOAPI_LOGGER_LEVEL", "LOGGER_LEVEL")
	c.Logger.Format = applyStringEnv(c.Logger.Format, "GOAPI_LOGGER_FORMAT", "LOGGER_FORMAT")

	c.JWT.Secret = applyStringEnv(c.JWT.Secret, "GOAPI_JWT_SECRET", "JWT_SECRET")
	c.JWT.ExpirySeconds = applyIntEnv(c.JWT.ExpirySeconds, "GOAPI_JWT_EXPIRY_SECONDS", "JWT_EXPIRY_SECONDS")
}

// DSN builds a PostgreSQL DSN string from the database config fields
// Format: host=... port=... user=... password=... dbname=... sslmode=...
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// Load reads config from YAML (when available), then applies environment overrides.
// If the config file does not exist, env-only configuration is supported.
// Returns a *Config on success, or an error for unreadable/invalid YAML files.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read config file %q: %w", path, err)
		}
	}

	var cfg Config
	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parse config file %q: %w", path, err)
		}
	}

	cfg.applyEnvOverrides()
	cfg.applyDefaults()

	return &cfg, nil
}
