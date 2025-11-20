package credentials

import "github.com/google/uuid"

func (uc *UseCase) UpdatePassword(userID uuid.UUID, newPassword string) (bool, error) {
	success, err := uc.repo.updatePassword(
		userID, newPassword)
	if err != nil {
		return false, err
	}
	return success, nil
}
