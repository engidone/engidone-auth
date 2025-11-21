package jwt

import "engidoneauth/internal/apperror"

func (uc *UseCase) RefreshToken(token, rfToken string) (*TokenInfo, error) {
	// Validate existing token
	err := uc.ValidateToken(token)

	if err != nil {
		return nil, err
	}

	userID, err := uc.repo.getUserByRefreshToken(rfToken)

	if err != nil {
		return nil, apperror.New(ErrInvalidRefreshToken, "Could not renew your session")
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
