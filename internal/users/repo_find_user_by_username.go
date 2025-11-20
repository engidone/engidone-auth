package users

func (r *RPCUserServiceRepository) findUserByUsername(username string) (*User, error) {

	for _, user := range r.users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, nil
}
