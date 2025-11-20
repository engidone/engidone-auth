package recovery

import "go.uber.org/fx"

var RecoveryModule = fx.Options(
	fx.Provide(NewSQLRepository),
	fx.Provide(NewUseCase),
)