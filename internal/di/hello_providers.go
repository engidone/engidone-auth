package di

import (
	"go.uber.org/fx"

	"engidone-auth/internal/hello/domain"
	"engidone-auth/internal/hello/infrastructure"
	"engidone-auth/internal/hello/usecase"
)

// HelloModule provides all hello service dependencies
var HelloModule = fx.Options(
	fx.Provide(
		NewHelloService,
		NewHelloUseCase,
	),
)

// NewHelloService provides a HelloService implementation
func NewHelloService() domain.HelloService {
	return infrastructure.NewHelloService()
}

// NewHelloUseCase provides a HelloUseCase implementation
func NewHelloUseCase(service domain.HelloService) domain.HelloUseCase {
	return usecase.NewHelloUseCase(service)
}