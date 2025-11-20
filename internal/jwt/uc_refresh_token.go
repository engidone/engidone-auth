package jwt

import "time"

func (uc *UseCase) RefreshToken(token string) (*TokenInfo, error) {
	// Validate existing token
	tokenInfo, err := uc.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Generate new token for the same user
	newToken, err := uc.GenerateToken(tokenInfo.UserID)
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		UserID:    tokenInfo.UserID,
		Token:     newToken,
		ExpiresAt: time.Now().Add(uc.tokenDuration),
		IssuedAt:  time.Now(),
	}, nil
}