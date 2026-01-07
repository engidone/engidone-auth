package jwt

import (
	"github.com/samber/oops"
)

// Error codes for JWT domain
const (
	CodeInvalidToken               = "JWT_INVALID_TOKEN"
	CodeTokenExpired               = "JWT_TOKEN_EXPIRED"
	CodeTokenGeneration            = "JWT_TOKEN_GENERATION_FAILED"
	CodeParsingToken              = "JWT_ERROR_PARSING_TOKEN"
	CodeGeneratingRefreshToken     = "JWT_ERROR_GENERATING_REFRESH_TOKEN"
	CodeInvalidRefreshToken        = "JWT_INVALID_REFRESH_TOKEN"
	CodeInsertOrUpdateRefreshToken = "JWT_ERROR_INSERT_OR_UPDATE_TOKEN"
)

// Domain-specific error builders
var (
	// Authentication errors
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

	// Compare with our predefined errors
	switch code {
	case CodeInvalidToken:
		return err == InvalidToken
	case CodeTokenExpired:
		return err == TokenExpired
	case CodeTokenGeneration:
		return err == TokenGenerationFailed
	case CodeParsingToken:
		return err == ParsingToken
	case CodeGeneratingRefreshToken:
		return err == GeneratingRefreshToken
	case CodeInvalidRefreshToken:
		return err == InvalidRefreshToken
	case CodeInsertOrUpdateRefreshToken:
		return err == InsertOrUpdateRefreshToken
	default:
		return false
	}
}

// ExtractJWTErrorInfo extracts JWT-specific error information
func ExtractJWTErrorInfo(err error) struct {
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

	// Check if this is a known JWT error
	switch err {
	case InvalidToken:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeInvalidToken,
			Message: "Invalid token",
			Details: map[string]interface{}{"domain": "jwt"},
		}
	case TokenExpired:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeTokenExpired,
			Message: "Token expired",
			Details: map[string]interface{}{"domain": "jwt"},
		}
	case TokenGenerationFailed:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeTokenGeneration,
			Message: "Token generation failed",
			Details: map[string]interface{}{"domain": "jwt"},
		}
	case ParsingToken:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeParsingToken,
			Message: "Error parsing token",
			Details: map[string]interface{}{"domain": "jwt"},
		}
	case GeneratingRefreshToken:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeGeneratingRefreshToken,
			Message: "Error generating refresh token",
			Details: map[string]interface{}{"domain": "jwt"},
		}
	case InvalidRefreshToken:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeInvalidRefreshToken,
			Message: "Invalid refresh token",
			Details: map[string]interface{}{"domain": "jwt"},
		}
	case InsertOrUpdateRefreshToken:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    CodeInsertOrUpdateRefreshToken,
			Message: "Error inserting or updating refresh token",
			Details: map[string]interface{}{"domain": "jwt"},
		}
	default:
		return struct {
			Code    string
			Message string
			Details map[string]interface{}
		}{
			Code:    "UNKNOWN_JWT_ERROR",
			Message: err.Error(),
			Details: map[string]interface{}{"domain": "jwt"},
		}
	}
}