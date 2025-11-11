package infrastructure

import (
	"fmt"
)

// HelloService implementa el servicio de saludo
type HelloService struct{}

// NewHelloService crea una nueva instancia del servicio de saludo
func NewHelloService() *HelloService {
	return &HelloService{}
}

// SayHello genera un saludo personalizado
func (s *HelloService) SayHello(name string) (string, error) {
	// Lógica simple para generar un saludo
	greeting := fmt.Sprintf("¡Hola, %s! Bienvenido al servicio de autenticación.", name)
	return greeting, nil
}