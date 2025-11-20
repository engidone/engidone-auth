package greet

import (
	"go.uber.org/fx"
)	

var GreetModule = fx.Options(
	fx.Provide(
		NewUseCase,
		NewGreetingRepository,
	),
)