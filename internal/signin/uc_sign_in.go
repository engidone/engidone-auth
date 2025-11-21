package signin

import (
	"engidoneauth/log"
	"fmt"
)

// Execute executes the authentication process
func (uc *UseCase) SingIn(credentials Credentials) (*Result, error) {
	// Validate credentials
	if err := uc.validateCredentials(credentials); err != nil {
		return nil, err
	}

	user, _ := uc.usersUC.GetUser(credentials.Username)

	// Verify user and password
	_, err := uc.credentialsUC.VerifyCredentials(user.ID, credentials.Password)

	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := uc.jwtUC.GenerateToken(user.ID)
	if err != nil {
		log.Error("Error generating token", log.String("user_id", user.ID), log.String("user_name", user.Username), log.Err(err))
		return nil, err
	}

	refreshToken, err := uc.jwtUC.GetRefreshToken()

	if err != nil {
		return nil, err
	}

	err = uc.jwtUC.SyncRefreshToken(user.ID, refreshToken)

	if err != nil {
		return nil, err
	}

	// Create authentication response
	result := &Result{
		Token:        token,
		RefreshToken: refreshToken,
	}

	fmt.Println("resut:: ", user)

	return result, nil
}
