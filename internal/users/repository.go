package users

type RPCUserServiceRepository struct {
	users []User
}

func NewRPCUserServiceRepository(users []User) *RPCUserServiceRepository {
	return &RPCUserServiceRepository{
		users: users,
	}
}
