package domain

import (
	"time"
)

// User representa la entidad de usuario en el dominio
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // No se serializa la contraseña
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Credentials representa las credenciales de autenticación
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse representa la respuesta de autenticación
type AuthResponse struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Token    string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// AuthError representa un error de autenticación
type AuthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AuthError) Error() string {
	return e.Message
}

// Constantes de errores de autenticación
const (
	ErrInvalidCredentials = "INVALID_CREDENTIALS"
	ErrUserNotFound       = "USER_NOT_FOUND"
	ErrUserDisabled       = "USER_DISABLED"
	ErrInvalidToken       = "INVALID_TOKEN"
)

// NewAuthError crea un nuevo error de autenticación
func NewAuthError(code, message string) *AuthError {
	return &AuthError{
		Code:    code,
		Message: message,
	}
}

// GetTokenExpiration extrae la fecha de expiración de un token (simulación)
func GetTokenExpiration(token string) time.Time {
	// Para esta implementación simple, devolvemos una expiración fija de 24 horas
	return time.Now().Add(24 * time.Hour)
}