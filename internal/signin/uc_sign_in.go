package signin

import (
	"engidoneauth/internal/apperror"
	"engidoneauth/log"
	"engidoneauth/util/timing"
)

// Execute executes the authentication process
func (uc *UseCase) SingIn(credentials Credentials) (*SigInResponse, error) {
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
		log.Error("Error generating token", log.Int32("user_id", user.ID), log.String("user_name", user.Username), log.Err(err))
		return nil, apperror.New(ErrInvalidToken, "Error generating token")
	}

	// Create authentication response
	response := &SigInResponse{
		Token:     token,
		ExpiresAt: timing.GetTokenExpiration(token),
	}

	return response, nil
}
