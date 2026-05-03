package utils

import (
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ParseOptionalUUIDToPGType converts an optional UUID string to pgtype.UUID.
// Returns an invalid (NULL) UUID when the input is empty or invalid.
// Sample input:  raw = "550e8400-e29b-41d4-a716-446655440000"
// Sample output: pgtype.UUID{Valid:true, Bytes:...}
func ParseOptionalUUIDToPGType(raw string) pgtype.UUID {
	value := strings.TrimSpace(raw)
	if value == "" {
		return pgtype.UUID{Valid: false}
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}

	return pgtype.UUID{Bytes: parsed, Valid: true}
}
