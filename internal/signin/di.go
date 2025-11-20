package signin

import "go.uber.org/fx"

var SignInModule = fx.Options(
	fx.Provide(
		NewUseCase,
	),
)
