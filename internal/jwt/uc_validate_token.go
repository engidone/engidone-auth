package jwt

import (
	"github.com/samber/oops"
	"github.com/golang-jwt/jwt/v5"
)

func (uc *UseCase) ValidateToken(tokenInput string) error {
	token, err := jwt.Parse(tokenInput, func(token *jwt.Token) (any, error) {
		return uc.certs.PublicKey, nil
	})

	if err != nil {
		return oops.With("token", tokenInput).Wrapf(ParsingToken, "Failed to parse JWT token: %v", err)
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	} else {
		return oops.With("token", tokenInput).Wrap(InvalidToken)
	}
}
