package jwt

func (uc *UseCase) SyncRefreshToken(userID, refreshToken string) error {
	_, err := uc.repo.syncRefreshToken(userID, refreshToken)
	return err
}
