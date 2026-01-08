package signin

import (
	"errors"

	"github.com/samber/oops"
)

// Error codes for signin domain
const (
	CodeInvalidCredentials = "SIGNIN_INVALID_CREDENTIALS"
	CodeMissingUsername    = "SIGNIN_MISSING_USERNAME"
	CodeMissingPassword    = "SIGNIN_MISSING_PASSWORD"
	CodeUsernameTooShort   = "SIGNIN_USERNAME_TOO_SHORT"
	CodePasswordTooShort   = "SIGNIN_PASSWORD_TOO_SHORT"
)

// Domain-specific error builders
var (
	// Authentication errors
	InvalidCredentials = oops.
				Code(CodeInvalidCredentials).
				With("domain", "signin").
				New("Invalid credentials")

	// Validation errors
	MissingUsername = oops.
			Code(CodeMissingUsername).
			With("domain", "signin").
			New("Username is required")

	MissingPassword = oops.
			Code(CodeMissingPassword).
			With("domain", "signin").
			New("Password is required")

	UsernameTooShort = oops.
				Code(CodeUsernameTooShort).
				With("domain", "signin").
				With("min_length", 3).
				New("Username must be at least 3 characters long")

	PasswordTooShort = oops.
				Code(CodePasswordTooShort).
				With("domain", "signin").
				With("min_length", 4).
				New("Password must be at least 4 characters long")
)

// IsErrorCode checks if an error has a specific signin error code
func IsErrorCode(err error, code string) bool {
	if err == nil {
		return false
	}

	// Extract error code from oops errors
	if oopsErr, ok := err.(*oops.OopsError); ok {
		return oopsErr.Code() == code
	}

	// Fallback to errors.Is for compatibility
	switch code {
	case CodeInvalidCredentials:
		return errors.Is(err, InvalidCredentials)
	case CodeMissingUsername:
		return errors.Is(err, MissingUsername)
	case CodeMissingPassword:
		return errors.Is(err, MissingPassword)
	case CodeUsernameTooShort:
		return errors.Is(err, UsernameTooShort)
	case CodePasswordTooShort:
		return errors.Is(err, PasswordTooShort)
	default:
		return false
	}
}
