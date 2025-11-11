package main

import (
	"go.uber.org/fx"
	"engidone-auth/internal/di"
)

func main() {
	app := fx.New(
		// Application-level providers (logger, config)
		di.LoggerModule,
		di.ConfigModule,

		// Domain-specific providers
		di.HelloModule,
		di.SigninModule,

		// gRPC transport providers
		di.GRPCModule,
	)

	app.Run()
}