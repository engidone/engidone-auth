package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"engidoneauth/internal/apperror"
	"fmt"
)


func (uc *UseCase) GenerateToken(userID string) (string, error) {
	// Generate a random simple token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", apperror.New(ErrInvalidToken, "Error generando token")
	}

	token := hex.EncodeToString(bytes)
	fullToken := fmt.Sprintf("Bearer %s", token)

	return fullToken, nil
}

