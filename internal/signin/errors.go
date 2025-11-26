package signin

import (
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

	// Compare with our predefined errors
	switch code {
	case CodeInvalidCredentials:
		return err == InvalidCredentials
	case CodeMissingUsername:
		return err == MissingUsername
	case CodeMissingPassword:
		return err == MissingPassword
	case CodeUsernameTooShort:
		return err == UsernameTooShort
	case CodePasswordTooShort:
		return err == PasswordTooShort
	default:
		return false
	}
}

// ExtractSigninErrorInfo extracts signin-specific error information
func ExtractSigninErrorInfo(err error) struct {
	Code    string
	Message string
	Details map[string]interface{}
} {
	if err == nil {
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{}
	}

	// Check if this is a known signin error
	switch err {
	case InvalidCredentials:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeInvalidCredentials,
			Message: "Invalid credentials",
			Details: map[string]interface{}{"domain": "signin"},
		}
	case MissingUsername:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeMissingUsername,
			Message: "Username is required",
			Details: map[string]interface{}{"domain": "signin"},
		}
	case MissingPassword:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeMissingPassword,
			Message: "Password is required",
			Details: map[string]interface{}{"domain": "signin"},
		}
	case UsernameTooShort:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeUsernameTooShort,
			Message: "Username must be at least 3 characters long",
			Details: map[string]interface{}{
				"domain":     "signin",
				"min_length": 3,
			},
		}
	case PasswordTooShort:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodePasswordTooShort,
			Message: "Password must be at least 4 characters long",
			Details: map[string]interface{}{
				"domain":     "signin",
				"min_length": 4,
			},
		}
	default:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    "UNKNOWN_SIGNIN_ERROR",
			Message: err.Error(),
			Details: map[string]interface{}{"domain": "signin"},
		}
	}
}