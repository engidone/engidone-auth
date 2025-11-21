package credentials

func (uc *UseCase) UpdatePassword(userID string, newPassword string) (bool, error) {
	success, err := uc.repo.updatePassword(
		userID, newPassword)
	if err != nil {
		return false, err
	}
	return success, nil
}
