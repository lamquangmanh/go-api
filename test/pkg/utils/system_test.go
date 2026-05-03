package utils_test

import (
	"testing"

	"github.com/google/uuid"

	"go-api/pkg/utils"
)

func TestGetTotalMemoryMB(t *testing.T) {
	if got := utils.GetTotalMemoryMB(); got < 0 {
		t.Fatalf("memory must be >= 0")
	}
}

func TestDetectHost(t *testing.T) {
	host, err := utils.DetectHost()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if host == nil {
		t.Fatalf("expected host info")
	}
	if host.CPU <= 0 {
		t.Fatalf("expected cpu > 0")
	}
}

func TestGenerateUUID(t *testing.T) {
	id := utils.GenerateUUID()
	if _, err := uuid.Parse(id); err != nil {
		t.Fatalf("expected valid uuid, got %q", id)
	}
}
