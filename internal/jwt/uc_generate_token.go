package jwt

import (
	"engidoneauth/internal/apperror"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (uc *UseCase) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(uc.certs.PrivateKey)

	if err != nil {
		return "", apperror.New(ErrInvalidToken, "Error signin token")
	}

	return tokenString, nil
}
