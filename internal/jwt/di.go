package jwt

import "go.uber.org/fx"

var RecoveryModule = fx.Options(
	fx.Provide(NewUseCase),
)
