package users

import (
	"go.uber.org/fx"
)

var UsersModule = fx.Options(
	fx.Provide(
		NewUseCase,
		provideUsers,
		NewRPCUserServiceRepository,
	),
	fx.Invoke(registerUsersLoader),
)

type UserConfig struct {
	Users []User `yaml:"users"`
}

func provideUsers() *[]User {
	return &[]User{}
}