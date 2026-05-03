package utils_test

import (
	"testing"

	"github.com/google/uuid"

	"go-api/pkg/utils"
)

func TestParseOptionalUUIDToPGType(t *testing.T) {
	if got := utils.ParseOptionalUUIDToPGType(""); got.Valid {
		t.Fatalf("expected empty input to return invalid pg uuid")
	}
	if got := utils.ParseOptionalUUIDToPGType("bad-uuid"); got.Valid {
		t.Fatalf("expected invalid input to return invalid pg uuid")
	}

	id := uuid.New().String()
	got := utils.ParseOptionalUUIDToPGType(id)
	if !got.Valid {
		t.Fatalf("expected valid uuid")
	}
	if uuid.UUID(got.Bytes).String() != id {
		t.Fatalf("unexpected parsed uuid: %s", uuid.UUID(got.Bytes).String())
	}
}
