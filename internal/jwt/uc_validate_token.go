package jwt

import (
	"engidoneauth/internal/apperror"
	"time"
)

func (uc *UseCase) ValidateToken(token string) (*TokenInfo, error) {
	if len(token) < 7 || token[:7] != "Bearer " {
		return nil, apperror.New(ErrInvalidToken, "Formato de token inválido")
	}

	// For this implementation, extract userID from token
	tokenData := token[7:] // Remove "Bearer "
	userID := uc.extractUserIDFromToken(tokenData)
	if userID == "" {
		return nil, apperror.New(ErrInvalidToken, "Token inválido o expirado")
	}

	return &TokenInfo{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(uc.tokenDuration),
		IssuedAt:  time.Now(),
	}, nil
}

func (uc *UseCase) extractUserIDFromToken(tokenData string) string {
	// In a real implementation, decode the token and extract userID
	// Here, we simulate by returning a fixed userID for demonstration
	return "user123"
}