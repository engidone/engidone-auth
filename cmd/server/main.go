package main

import (
	"engidoneauth/internal/config"
	"engidoneauth/internal/greet"
	"engidoneauth/internal/jwt"
	"engidoneauth/internal/recovery"
	"engidoneauth/internal/server"
	"engidoneauth/internal/signin"
	"engidoneauth/internal/users"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		config.ConfigModule,
		server.ServerModule,
		jwt.RecoveryModule,
		greet.GreetModule,
		users.UsersModule,
		signin.SignInModule,
		recovery.RecoveryModule,
	)

	app.Run()
}
