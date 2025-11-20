package config

import "go.uber.org/fx"

var ConfigModule = fx.Options(
	fx.Provide(newConfigPaths),
	fx.Provide(newAppConfig),
)
