package users

import (
	"errors"

	"github.com/samber/oops"
)

// Error codes for users domain
const (
	CodeUserNotFound      = "USERS_USER_NOT_FOUND"
	CodeUserAlreadyExists = "USERS_USER_ALREADY_EXISTS"
	CodeUserInvalid       = "USERS_USER_INVALID"
)

// Domain-specific error builders
var (
	// User domain errors
	UserNotFound = oops.
			Code(CodeUserNotFound).
			With("domain", "users").
			New("User not found")

	UserAlreadyExists = oops.
				Code(CodeUserAlreadyExists).
				With("domain", "users").
				New("User already exists")

	UserInvalid = oops.
			Code(CodeUserInvalid).
			With("domain", "users").
			New("Invalid user data")
)

// IsErrorCode checks if an error has a specific users error code
func IsErrorCode(err error, code string) bool {
	if err == nil {
		return false
	}

	// Extract error code from oops errors
	if oopsErr, ok := err.(*oops.OopsError); ok {
		return oopsErr.Code() == code
	}

	// Extract error code from oops errors
	if oopsErr, ok := err.(*oops.OopsError); ok {
		return oopsErr.Code() == code
	}

	// Fallback to errors.Is for compatibility
	switch code {
	case CodeUserNotFound:
		return errors.Is(err, UserNotFound)
	case CodeUserAlreadyExists:
		return errors.Is(err, UserAlreadyExists)
	case CodeUserInvalid:
		return errors.Is(err, UserInvalid)
	default:
		return false
	}
}