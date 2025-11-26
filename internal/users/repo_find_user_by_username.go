package users

import (
	"engidoneauth/util/collection"
	"github.com/samber/oops"
)

func (r *RPCUserServiceRepository) findUserByUsername(username string) (*User, error) {

	user, found := collection.Find(r.users, func(user User) bool {
		return user.Username == username
	})

	if !found {
		return nil, oops.With("username", username).Wrap(UserNotFound)
	}

	return &user, nil
}
