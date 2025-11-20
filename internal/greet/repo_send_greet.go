package greet

import (
	"fmt"
)

// SayHello generates a personalized greeting
func (s *greetingRepository) sendGreeting(name string) (string, error) {
	// Simple logic to generate a greeting
	greeting := fmt.Sprintf("¡Hola, %s! Bienvenido al servicio de autenticación.", name)
	return greeting, nil
}