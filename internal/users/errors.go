package users

import (
	"github.com/samber/oops"
)

// Error codes for users domain
const (
	CodeUserNotFound     = "USERS_USER_NOT_FOUND"
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

	// Compare with our predefined errors
	switch code {
	case CodeUserNotFound:
		return err == UserNotFound
	case CodeUserAlreadyExists:
		return err == UserAlreadyExists
	case CodeUserInvalid:
		return err == UserInvalid
	default:
		return false
	}
}

// ExtractUserErrorInfo extracts user-specific error information
func ExtractUserErrorInfo(err error) struct {
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

	// Check if this is a known users error
	switch err {
	case UserNotFound:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeUserNotFound,
			Message: "User not found",
			Details: map[string]interface{}{"domain": "users"},
		}
	case UserAlreadyExists:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeUserAlreadyExists,
			Message: "User already exists",
			Details: map[string]interface{}{"domain": "users"},
		}
	case UserInvalid:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeUserInvalid,
			Message: "Invalid user data",
			Details: map[string]interface{}{"domain": "users"},
		}
	default:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    "UNKNOWN_USER_ERROR",
			Message: err.Error(),
			Details: map[string]interface{}{"domain": "users"},
		}
	}
}