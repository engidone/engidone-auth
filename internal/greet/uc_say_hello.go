package greet

import (
	"engidoneauth/internal/apperror"
)

// Execute executes the hello use case
func (uc *UseCase) SayHello(name string) (*HelloResponse, error) {
	if name == "" {
		return &HelloResponse{
			Message: "Por favor, proporciona un nombre",
			Success: false,
		}, apperror.New(ErrEmptyName,"nombre vacío")	
	}

	message, err := uc.repository.sendGreeting(name)
	if err != nil {
		return &HelloResponse{
			Message: "Error al generar saludo",
			Success: false,
		}, err
	}

	return &HelloResponse{
		Message: message,
		Success: true,
	}, nil
}