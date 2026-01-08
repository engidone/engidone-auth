package recovery

import (
	"errors"

	"github.com/samber/oops"
)

// Error codes for recovery domain
const (
	CodeRecoveryCodeNotFound = "RECOVERY_CODE_NOT_FOUND"
	CodeInternalError        = "RECOVERY_INTERNAL_ERROR"
	CodeRecoveryCodeInvalid  = "RECOVERY_CODE_INVALID"
)

// Domain-specific error builders
var (
	// Recovery code errors
	RecoveryCodeNotFound = oops.
				Code(CodeRecoveryCodeNotFound).
				With("domain", "recovery").
				New("Recovery code not found")

	// System errors
	InternalError = oops.
			Code(CodeInternalError).
			With("domain", "recovery").
			New("Internal error")

	RecoveryCodeInvalid = oops.
				Code(CodeInternalError).
				With("domain", "recovery").
				New("Internal error")
)

// IsErrorCode checks if an error has a specific recovery error code
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
	case CodeRecoveryCodeNotFound:
		return errors.Is(err, RecoveryCodeNotFound)
	case CodeInternalError:
		return errors.Is(err, InternalError)
	case CodeRecoveryCodeInvalid:
		return errors.Is(err, RecoveryCodeInvalid)
	default:
		return false
	}
}
