package logger

import (
	"log/slog"
	"os"
	"strings"

	"go-api/internal/config"
)

// New creates a structured logger from runtime config.
// - level: debug | info | warn | error
// - format: json (default) | console/text
// Invalid or empty values fall back to safe defaults.
func New(cfg config.LoggerConfig) (*slog.Logger, error) {
	var level slog.Level
	// Normalize and map string log level to slog level constants.
	switch strings.ToLower(strings.TrimSpace(cfg.Level)) {
	case "debug":
		level = slog.LevelDebug
	case "info", "":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}
	// Choose output format. JSON is preferred for machine parsing.
	format := strings.ToLower(strings.TrimSpace(cfg.Format))
	if format == "console" || format == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, opts)), nil
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, opts)), nil
}
