package credentials

import (
	"engidoneauth/internal/db"

	"github.com/google/uuid"
)

type repository interface {
	updatePassword(userID uuid.UUID, newPassword string) (bool, error)
	findCredential(userID uuid.UUID, password string) (*db.Credential, error)
}

type UseCase struct {
	repo repository
}

func NewUseCase(credentialsRepository repository) *UseCase {
	return &UseCase{
		repo: credentialsRepository,
	}
}
