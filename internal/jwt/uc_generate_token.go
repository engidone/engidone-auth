package jwt

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/oops"
)

func (uc *UseCase) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(uc.certs.PrivateKey)

	if err != nil {
		return "", oops.
			With("user_id", userID).
			Wrapf(GeneratingRefreshToken, "Failed to sign JWT token: %v", err)
	}

	return tokenString, nil
}
