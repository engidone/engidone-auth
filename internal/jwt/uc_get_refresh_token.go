package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"engidoneauth/internal/apperror"
)

func (uc *UseCase) GetRefreshToken() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)

	if err != nil {
		return "", apperror.New(ErrGeneratingRefreshToken, err.Error())
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
