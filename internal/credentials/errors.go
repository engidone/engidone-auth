package credentials

import (
	"errors"

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

	// Extract error code from oops errors
	if oopsErr, ok := err.(*oops.OopsError); ok {
		return oopsErr.Code() == code
	}

	switch code {
	case CodeInvalidCredentials:
		return errors.Is(err, InvalidCredentials)
	case CodeUpdatePasswordFailed:
		return errors.Is(err, UpdatePasswordFailed)
	case CodeInternalError:
		return errors.Is(err, InternalError)
	case CodeUserNotFound:
		return errors.Is(err, UserNotFound)
	default:
		return false
	}

}
