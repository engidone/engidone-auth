package signin

func (uc *UseCase) validateCredentials(credentials Credentials) error {
	if credentials.Username == "" {
		return MissingUsername
	}

	if credentials.Password == "" {
		return MissingPassword
	}

	if len(credentials.Username) < 3 {
		return UsernameTooShort
	}

	if len(credentials.Password) < 4 {
		return PasswordTooShort
	}

	return nil
}