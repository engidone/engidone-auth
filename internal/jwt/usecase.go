package jwt

import (
	"engidoneauth/internal/db"
	"time"
)

type repository interface {
	syncRefreshToken(userID, refreshToken string) (*db.RefreshToken, error)
	getUserByRefreshToken(refreshToken string) (string, error)
}

type UseCase struct {
	tokenDuration time.Duration
	certs         Certs
	repo repository
}

func NewUseCase(tokenDuration time.Duration, certs Certs, repo repository) *UseCase {
	return &UseCase{tokenDuration, certs, repo}
}
