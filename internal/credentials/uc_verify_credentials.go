package credentials

import (
	"engidoneauth/internal/apperror"
	"engidoneauth/util/crypto"

	"github.com/google/uuid"
)

func (uc *UseCase) VerifyCredentials(userID uuid.UUID, password string) (bool, error) {

	// Hash the provided password
	hashedPassword := crypto.HashPassword(password)

	storeduser, err := uc.repo.findCredential(userID, hashedPassword)

	if err != nil {
		return false, err
	}

	if storeduser == nil || storeduser.UserID != userID {
		return false, apperror.New(ErrInvalidCredentials, "Invalid credentials")
	}

	return true, nil
}
