package users

func (uc *UseCase) GetUser(username string) (*User, error) {
	user, err := uc.repo.findUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}