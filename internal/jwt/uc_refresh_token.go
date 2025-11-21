package jwt

func (uc *UseCase) RefreshToken(token, rfToken, userID string) (*TokenInfo, error) {
	// Validate existing token
	err := uc.ValidateToken(token)

	if err != nil {
		return nil, err
	}

	// Generate new token for the same user
	newToken, err := uc.GenerateToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.GetRefreshToken()

	if err != nil {
		return nil, err
	}

	_, err = uc.repo.syncRefreshToken(userID, refreshToken)

	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		Token:        newToken,
		RefreshToken: refreshToken,
	}, nil
}
