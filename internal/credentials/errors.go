package credentials

import (
	"github.com/samber/oops"
)

// Error codes for credentials domain
const (
	CodeInvalidCredentials   = "CREDENTIALS_INVALID"
	CodeUpdatePasswordFailed = "CREDENTIALS_UPDATE_FAILED"
	CodeInternalError        = "CREDENTIALS_INTERNAL_ERROR"
	CodeUserNotFound         = "CREDENTIALS_USER_NOT_FOUND"
)

// Domain-specific error builders
var (
	// Authentication errors
	InvalidCredentials = oops.
		Code(CodeInvalidCredentials).
		With("domain", "credentials").
		New("Invalid credentials")

	// Password management errors
	UpdatePasswordFailed = oops.
		Code(CodeUpdatePasswordFailed).
		With("domain", "credentials").
		New("Update password failed")

	// System errors
	InternalError = oops.
		Code(CodeInternalError).
		With("domain", "credentials").
		New("Internal error")

	// User related errors
	UserNotFound = oops.
		Code(CodeUserNotFound).
		With("domain", "credentials").
		New("User not found")
)

// IsErrorCode checks if an error has a specific credentials error code
func IsErrorCode(err error, code string) bool {
	if err == nil {
		return false
	}

	// Compare with our predefined errors
	switch code {
	case CodeInvalidCredentials:
		return err == InvalidCredentials
	case CodeUpdatePasswordFailed:
		return err == UpdatePasswordFailed
	case CodeInternalError:
		return err == InternalError
	case CodeUserNotFound:
		return err == UserNotFound
	default:
		return false
	}
}

// ExtractCredentialsErrorInfo extracts credentials-specific error information
func ExtractCredentialsErrorInfo(err error) struct {
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

	// Check if this is a known credentials error
	switch err {
	case InvalidCredentials:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeInvalidCredentials,
			Message: "Invalid credentials",
			Details: map[string]interface{}{"domain": "credentials"},
		}
	case UpdatePasswordFailed:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeUpdatePasswordFailed,
			Message: "Update password failed",
			Details: map[string]interface{}{"domain": "credentials"},
		}
	case InternalError:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeInternalError,
			Message: "Internal error",
			Details: map[string]interface{}{"domain": "credentials"},
		}
	case UserNotFound:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeUserNotFound,
			Message: "User not found",
			Details: map[string]interface{}{"domain": "credentials"},
		}
	default:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    "UNKNOWN_CREDENTIALS_ERROR",
			Message: err.Error(),
			Details: map[string]interface{}{"domain": "credentials"},
		}
	}
}