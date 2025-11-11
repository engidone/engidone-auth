package usecase

import (
	"errors"
	"engidone-auth/internal/hello/domain"
)

// HelloUseCase implementa el caso de uso para el servicio de saludo
type HelloUseCase struct {
	helloService domain.HelloService
}

// NewHelloUseCase crea una nueva instancia del caso de uso de saludo
func NewHelloUseCase(helloService domain.HelloService) *HelloUseCase {
	return &HelloUseCase{
		helloService: helloService,
	}
}

// Execute ejecuta el caso de uso de saludo
func (uc *HelloUseCase) Execute(name string) (*domain.HelloResponse, error) {
	if name == "" {
		return &domain.HelloResponse{
			Message: "Por favor, proporciona un nombre",
			Success: false,
		}, errors.New("nombre vacío")
	}

	message, err := uc.helloService.SayHello(name)
	if err != nil {
		return &domain.HelloResponse{
			Message: "Error al generar saludo",
			Success: false,
		}, err
	}

	return &domain.HelloResponse{
		Message: message,
		Success: true,
	}, nil
}