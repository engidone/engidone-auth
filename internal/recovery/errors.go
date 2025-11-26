package recovery

import (
	"github.com/samber/oops"
)

// Error codes for recovery domain
const (
	CodeRecoveryCodeNotFound = "RECOVERY_CODE_NOT_FOUND"
	CodeDatabaseError        = "RECOVERY_DATABASE_ERROR"
	CodeInternalError        = "RECOVERY_INTERNAL_ERROR"
	CodeTimeoutError         = "RECOVERY_TIMEOUT_ERROR"
)

// Domain-specific error builders
var (
	// Recovery code errors
	RecoveryCodeNotFound = oops.
		Code(CodeRecoveryCodeNotFound).
		With("domain", "recovery").
		New("Recovery code not found")

	// Database errors
	DatabaseError = oops.
		Code(CodeDatabaseError).
		With("domain", "recovery").
		New("Database operation failed")

	// System errors
	InternalError = oops.
		Code(CodeInternalError).
		With("domain", "recovery").
		New("Internal error")

	// Timeout errors
	TimeoutError = oops.
		Code(CodeTimeoutError).
		With("domain", "recovery").
		New("Operation timed out")
)

// IsErrorCode checks if an error has a specific recovery error code
func IsErrorCode(err error, code string) bool {
	if err == nil {
		return false
	}

	// Compare with our predefined errors
	switch code {
	case CodeRecoveryCodeNotFound:
		return err == RecoveryCodeNotFound
	case CodeDatabaseError:
		return err == DatabaseError
	case CodeInternalError:
		return err == InternalError
	case CodeTimeoutError:
		return err == TimeoutError
	default:
		return false
	}
}

// ExtractRecoveryErrorInfo extracts recovery-specific error information
func ExtractRecoveryErrorInfo(err error) struct {
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

	// Check if this is a known recovery error
	switch err {
	case RecoveryCodeNotFound:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeRecoveryCodeNotFound,
			Message: "Recovery code not found",
			Details: map[string]interface{}{"domain": "recovery"},
		}
	case DatabaseError:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeDatabaseError,
			Message: "Database operation failed",
			Details: map[string]interface{}{"domain": "recovery"},
		}
	case InternalError:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeInternalError,
			Message: "Internal error",
			Details: map[string]interface{}{"domain": "recovery"},
		}
	case TimeoutError:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeTimeoutError,
			Message: "Operation timed out",
			Details: map[string]interface{}{"domain": "recovery"},
		}
	default:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    "UNKNOWN_RECOVERY_ERROR",
			Message: err.Error(),
			Details: map[string]interface{}{"domain": "recovery"},
		}
	}
}