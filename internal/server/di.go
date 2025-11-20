package server

import "go.uber.org/fx"

var ServerModule = fx.Options(
	fx.Provide(newGRPCServer),
	fx.Provide(newDB),
	fx.Provide(newQueries),
	fx.Invoke(registerDBConnectionLoader),
)
