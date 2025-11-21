package jwt

import (
	"engidoneauth/internal/db"
	"time"
)

type repository interface {
	existsRefreshToken(userID, refreshToken string) (bool, error)
	syncRefreshToken(userID, refreshToken string) (*db.RefreshToken, error)
}

type UseCase struct {
	tokenDuration time.Duration
	certs         Certs
	repo repository
}

func NewUseCase(tokenDuration time.Duration, certs Certs, repo repository) *UseCase {
	return &UseCase{tokenDuration, certs, repo}
}
