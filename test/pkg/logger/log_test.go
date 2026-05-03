package logger_test

import (
	"testing"

	"go-api/pkg/logger"
)

func TestInitLoggerAndPublicAPIs(t *testing.T) {
	logger.InitLogger()
	logger.Info("info: %s", "ok")
	logger.Warn("warn: %s", "ok")
	logger.Error("error: %s", "ok")
	logger.Debug("debug: %s", "ok")
}
