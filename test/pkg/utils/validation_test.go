package utils_test

import (
	"testing"

	"go-api/pkg/utils"
)

func TestValidateUsername(t *testing.T) {
	if !utils.ValidateUsername("john.doe_1") {
		t.Fatalf("expected valid username")
	}
	if utils.ValidateUsername("ab") {
		t.Fatalf("expected too-short username to be invalid")
	}
	if utils.ValidateUsername("bad user") {
		t.Fatalf("expected username with spaces to be invalid")
	}
}

func TestValidateEmail(t *testing.T) {
	if !utils.ValidateEmail("john@example.com") {
		t.Fatalf("expected valid email")
	}
	if utils.ValidateEmail("not-an-email") {
		t.Fatalf("expected invalid email")
	}
}

func TestValidateRoleNameAndMaxLength(t *testing.T) {
	if !utils.ValidateRoleName("Admin") {
		t.Fatalf("expected valid role name")
	}
	if utils.ValidateRoleName("a") {
		t.Fatalf("expected too-short role name to be invalid")
	}
	if !utils.ValidateMaxLength("  abc  ", 3) {
		t.Fatalf("expected trimmed length check to pass")
	}
	if utils.ValidateMaxLength("abcd", 3) {
		t.Fatalf("expected max length validation to fail")
	}
}
