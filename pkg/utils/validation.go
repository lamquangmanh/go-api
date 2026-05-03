package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// sharedValidate is a reusable validator instance for generic input checks.
var sharedValidate = validator.New()

// ValidateUsername checks username format and length.
// This function is intentionally generic and does not return business-specific
// error messages; callers should map false to domain-specific messages.
// Sample input:  "john.doe_1"
// Sample output: true
func ValidateUsername(username string) bool {
	value := strings.TrimSpace(username)
	if value == "" {
		return false
	}

	if err := sharedValidate.Var(value, "min=3,max=50"); err != nil {
		return false
	}

	for _, ch := range value {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '_' || ch == '.' || ch == '-') {
			return false
		}
	}

	return true
}

// ValidateEmail checks email format and max length.
// This function is intentionally generic and does not return business-specific
// error messages; callers should map false to domain-specific messages.
// Sample input:  "john@example.com"
// Sample output: true
func ValidateEmail(email string) bool {
	value := strings.TrimSpace(email)
	if value == "" {
		return false
	}

	if err := sharedValidate.Var(value, "email,max=254"); err != nil {
		return false
	}

	return true
}

// ValidateRoleName checks role name length constraints.
// This function is intentionally generic and does not return business-specific
// error messages; callers should map false to domain-specific messages.
// Sample input:  "Admin"
// Sample output: true
func ValidateRoleName(name string) bool {
	value := strings.TrimSpace(name)
	if value == "" {
		return false
	}

	if err := sharedValidate.Var(value, "min=2,max=100"); err != nil {
		return false
	}

	return true
}

// ValidateMaxLength checks whether a trimmed string length is <= max.
// Sample input:  value="  abc  ", max=3
// Sample output: true
func ValidateMaxLength(value string, max int) bool {
	return len(strings.TrimSpace(value)) <= max
}
