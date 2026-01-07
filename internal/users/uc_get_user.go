package users

import (
	"github.com/samber/oops"
)

func (uc *UseCase) GetUser(username string) (*User, error) {
	user, err := uc.repo.findUserByUsername(username)
	if err != nil {
		return nil, oops.Wrapf(err, "Failed to search user by username: %s", username)
	}
	return user, nil
}