package signin

import "engidoneauth/internal/apperror"

func (uc *UseCase) validateCredentials(credentials Credentials) error {
	if credentials.Username == "" {
		return apperror.New(ErrInvalidCredentials, "The username is required")
	}

	if credentials.Password == "" {
		return apperror.New(ErrInvalidCredentials, "The password is required")
	}

	if len(credentials.Username) < 3 {
		return apperror.New(ErrInvalidCredentials, "The username must be at least 3 characters long")
	}

	if len(credentials.Password) < 4 {
		return apperror.New(ErrInvalidCredentials, "The password must be at least 4 characters long")
	}

	return nil
}