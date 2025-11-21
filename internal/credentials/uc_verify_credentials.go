package credentials

import (
	"engidoneauth/internal/apperror"
	"engidoneauth/util/crypto"
)

func (uc *UseCase) VerifyCredentials(userID string, password string) (bool, error) {

	// Hash the provided password
	hashedPassword := crypto.HashPassword(password)

	storedUser, err := uc.repo.findCredential(userID, hashedPassword)

	if err != nil {
		return false, err
	}

	if storedUser == nil || storedUser.UserID.String() != userID {
		return false, apperror.New(ErrInvalidCredentials, "Invalid credentials")
	}

	return true, nil
}
