package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/samber/oops"
)

func (uc *UseCase) GetRefreshToken() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)

	if err != nil {
		return "", oops.Wrapf(GeneratingRefreshToken, "Failed to generate random bytes: %v", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
