package test

import (
	"os"
	"path/filepath"
	"testing"

	"go-api/internal/config"
)

func TestLoad_ConfigFileWithEnvOverrides(t *testing.T) {
	// This test verifies that environment variables override values from YAML config.
	// This behavior is required for Kubernetes deployments where secrets are injected via env.
	//
	// Arrange:
	// - Set environment variables for selected fields.
	// - Write a temporary YAML file with different values.
	//
	// Act:
	// - Load configuration from the temporary file.
	//
	// Assert:
	// - Env values must win over file values for the same fields.
	t.Setenv("GOAPI_GRPC_PORT", "9000")
	t.Setenv("GOAPI_DB_HOST", "db-from-env")
	t.Setenv("GOAPI_DB_PASSWORD", "secret-from-env")

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yml")

	content := []byte(`grpc:
  port: 50051
http:
  port: 8080
database:
  host: db-from-file
  port: 5432
  user: postgres
  password: from-file
  name: go_api
  ssl_mode: disable
logger:
  level: info
  format: json
`)
	if err := os.WriteFile(cfgPath, content, 0o600); err != nil {
		// Fail fast when test fixtures cannot be prepared.
		t.Fatalf("write temp config: %v", err)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		// Config loading is the behavior under test; any error is a hard failure.
		t.Fatalf("load config: %v", err)
	}

	// GOAPI_GRPC_PORT must override YAML grpc.port.
	if cfg.GRPC.Port != 9000 {
		t.Fatalf("expected grpc port 9000 from env, got %d", cfg.GRPC.Port)
	}
	// GOAPI_DB_HOST must override YAML database.host.
	if cfg.Database.Host != "db-from-env" {
		t.Fatalf("expected db host from env, got %q", cfg.Database.Host)
	}
	// GOAPI_DB_PASSWORD must override YAML database.password.
	if cfg.Database.Password != "secret-from-env" {
		t.Fatalf("expected db password from env, got %q", cfg.Database.Password)
	}
}

func TestLoad_MissingFileUsesEnvAndDefaults(t *testing.T) {
	// This test verifies env-only startup mode (no YAML file on disk),
	// which is common in containers where runtime config comes from secrets.
	//
	// Arrange:
	// - Provide essential database values via environment variables only.
	// - Use a path that does not contain a config file.
	//
	// Act:
	// - Load configuration from the missing file path.
	//
	// Assert:
	// - Loader must not fail.
	// - Env values must be applied.
	// - Built-in defaults must be used for fields not provided in env.
	t.Setenv("GOAPI_DB_HOST", "db.svc.cluster.local")
	t.Setenv("GOAPI_DB_PORT", "5432")
	t.Setenv("GOAPI_DB_USER", "postgres")
	t.Setenv("GOAPI_DB_PASSWORD", "super-secret")
	t.Setenv("GOAPI_DB_NAME", "go_api")

	nonExistentPath := filepath.Join(t.TempDir(), "missing-config.yml")

	cfg, err := config.Load(nonExistentPath)
	if err != nil {
		// Missing file is expected in env-only mode; loader should still succeed.
		t.Fatalf("load config from env-only mode: %v", err)
	}

	// Default gRPC port should be applied when not specified via env.
	if cfg.GRPC.Port != 50051 {
		t.Fatalf("expected default grpc port 50051, got %d", cfg.GRPC.Port)
	}
	// Default HTTP port should be applied when not specified via env.
	if cfg.HTTP.Port != 8080 {
		t.Fatalf("expected default http port 8080, got %d", cfg.HTTP.Port)
	}
	// DB host should come from GOAPI_DB_HOST.
	if cfg.Database.Host != "db.svc.cluster.local" {
		t.Fatalf("expected db host from env, got %q", cfg.Database.Host)
	}
	// DB password should come from GOAPI_DB_PASSWORD.
	if cfg.Database.Password != "super-secret" {
		t.Fatalf("expected db password from env, got %q", cfg.Database.Password)
	}
	// Logger defaults should be applied when logger env vars are absent.
	if cfg.Logger.Level != "info" {
		t.Fatalf("expected default logger level info, got %q", cfg.Logger.Level)
	}
	if cfg.Logger.Format != "json" {
		t.Fatalf("expected default logger format json, got %q", cfg.Logger.Format)
	}
}
