package credentials

import (
	"engidoneauth/util/crypto"
	"github.com/samber/oops"
)

func (uc *UseCase) VerifyCredentials(userID string, password string) (bool, error) {
	// Hash the provided password
	hashedPassword := crypto.HashPassword(password)

	storedUser, err := uc.repo.findCredential(userID, hashedPassword)
	if err != nil {
		return false, err
	}

	if storedUser == nil || storedUser.UserID.String() != userID {
		return false, oops.
			With("user_id", userID).
			Wrap(InvalidCredentials)
	}

	return true, nil
}
