package credentials

import "go.uber.org/fx"

var CredentialsModule = fx.Options(
	fx.Provide(NewSQLRepository),
	fx.Provide(NewUseCase),
)