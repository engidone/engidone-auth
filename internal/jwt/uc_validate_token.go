package jwt

import (
	"engidoneauth/internal/apperror"

	"github.com/golang-jwt/jwt/v5"
)

func (uc *UseCase) ValidateToken(tokenInput string) error {
	token, err := jwt.Parse(tokenInput, func(token *jwt.Token) (any, error) {
		return uc.certs.PublicKey, nil
	})

	if err != nil {
		return apperror.New(ErrParsingToken, err.Error())
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	} else {
		return apperror.New(ErrInvalidToken, "Token invalid or expired")
	}
}
