package credentials

import (
	"engidoneauth/internal/db"
)

type repository interface {
	updatePassword(userID string, newPassword string) (bool, error)
	findCredential(userID string, password string) (*db.Credential, error)
}

type UseCase struct {
	repo repository
}

func NewUseCase(credentialsRepository repository) *UseCase {
	return &UseCase{
		repo: credentialsRepository,
	}
}
