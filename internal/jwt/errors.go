package jwt

import (
	"errors"

	"github.com/samber/oops"
)

// Error codes for JWT domain
const (
	CodeNotFoundRefreshToken       = "JWT_REFRESH_TOKEN_NOT_FOUND"
	CodeInvalidToken               = "JWT_INVALID_TOKEN"
	CodeTokenExpired               = "JWT_TOKEN_EXPIRED"
	CodeTokenGeneration            = "JWT_TOKEN_GENERATION_FAILED"
	CodeParsingToken               = "JWT_ERROR_PARSING_TOKEN"
	CodeGeneratingRefreshToken     = "JWT_ERROR_GENERATING_REFRESH_TOKEN"
	CodeInvalidRefreshToken        = "JWT_INVALID_REFRESH_TOKEN"
	CodeInsertOrUpdateRefreshToken = "JWT_ERROR_INSERT_OR_UPDATE_TOKEN"
)

// Domain-specific error builders
var (
	// Authentication errors

	RefreshTokenNotFound = oops.
				Code(CodeNotFoundRefreshToken).
				With("domain", "jwt").
				New("Refresh token not found")

	InvalidToken = oops.
			Code(CodeInvalidToken).
			With("domain", "jwt").
			New("Invalid token")

	TokenExpired = oops.
			Code(CodeTokenExpired).
			With("domain", "jwt").
			New("Token expired")

	TokenGenerationFailed = oops.
				Code(CodeTokenGeneration).
				With("domain", "jwt").
				New("Token generation failed")

	// Token parsing errors
	ParsingToken = oops.
			Code(CodeParsingToken).
			With("domain", "jwt").
			New("Error parsing token")

	// Refresh token errors
	GeneratingRefreshToken = oops.
				Code(CodeGeneratingRefreshToken).
				With("domain", "jwt").
				New("Error generating refresh token")

	InvalidRefreshToken = oops.
				Code(CodeInvalidRefreshToken).
				With("domain", "jwt").
				New("Invalid refresh token")

	InsertOrUpdateRefreshToken = oops.
					Code(CodeInsertOrUpdateRefreshToken).
					With("domain", "jwt").
					New("Error inserting or updating refresh token")
)

// IsErrorCode checks if an error has a specific JWT error code
func IsErrorCode(err error, code string) bool {
	if err == nil {
		return false
	}

	// Try to extract oops.Error and compare codes
	if oopsErr, ok := oops.AsOops(err); ok {
		return oopsErr.Code() == code
	}

	// Fallback to errors.Is for compatibility
	switch code {
	case CodeInvalidToken:
		return errors.Is(err, InvalidToken)
	case CodeTokenExpired:
		return errors.Is(err, TokenExpired)
	case CodeTokenGeneration:
		return errors.Is(err, TokenGenerationFailed)
	case CodeParsingToken:
		return errors.Is(err, ParsingToken)
	case CodeGeneratingRefreshToken:
		return errors.Is(err, GeneratingRefreshToken)
	case CodeInvalidRefreshToken:
		return errors.Is(err, InvalidRefreshToken)
	case CodeInsertOrUpdateRefreshToken:
		return errors.Is(err, InsertOrUpdateRefreshToken)
	default:
		return false
	}
}
