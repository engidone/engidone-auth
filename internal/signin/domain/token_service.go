package domain

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// TokenService define la interfaz para operaciones con tokens
type TokenService interface {
	// GenerateToken genera un nuevo token para un usuario
	GenerateToken(userID string) (string, error)

	// ValidateToken valida un token y extrae el userID
	ValidateToken(token string) (*TokenInfo, error)

	// RefreshToken genera un nuevo token refrescando uno existente
	RefreshToken(token string) (*TokenInfo, error)
}

// TokenInfo contiene la información extraída de un token
type TokenInfo struct {
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt time.Time `json:"issued_at"`
}

// JWTTokenService implementa TokenService con tokens simples (no JWT real para simplificar)
type JWTTokenService struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTTokenService crea una nueva instancia del servicio de tokens
func NewJWTTokenService() TokenService {
	return &JWTTokenService{
		secretKey:     "your-secret-key-change-in-production",
		tokenDuration: 24 * time.Hour, // 24 horas
	}
}

// GenerateToken genera un nuevo token para un usuario
func (s *JWTTokenService) GenerateToken(userID string) (string, error) {
	// Generar un token aleatorio simple
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", NewAuthError(ErrInvalidToken, "Error generando token")
	}

	token := hex.EncodeToString(bytes)
	fullToken := fmt.Sprintf("Bearer %s", token)

	return fullToken, nil
}

// ValidateToken valida un token y extrae la información
func (s *JWTTokenService) ValidateToken(token string) (*TokenInfo, error) {
	if len(token) < 7 || token[:7] != "Bearer " {
		return nil, NewAuthError(ErrInvalidToken, "Formato de token inválido")
	}

	// Para esta implementación simple, extraemos el userID del token
	// En una implementación real, usarías JWT y validarías la firma
	tokenData := token[7:] // Remover "Bearer "

	// Simular extracción de userID (en JWT real sería del payload)
	userID := s.extractUserIDFromToken(tokenData)
	if userID == "" {
		return nil, NewAuthError(ErrInvalidToken, "Token inválido o expirado")
	}

	return &TokenInfo{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.tokenDuration),
		IssuedAt:  time.Now(),
	}, nil
}

// RefreshToken genera un nuevo token refrescando uno existente
func (s *JWTTokenService) RefreshToken(token string) (*TokenInfo, error) {
	// Validar token existente
	tokenInfo, err := s.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Generar nuevo token para el mismo usuario
	newToken, err := s.GenerateToken(tokenInfo.UserID)
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		UserID:    tokenInfo.UserID,
		Token:     newToken,
		ExpiresAt: time.Now().Add(s.tokenDuration),
		IssuedAt:  time.Now(),
	}, nil
}

// extractUserIDFromToken extrae el userID de un token (implementación simple)
func (s *JWTTokenService) extractUserIDFromToken(token string) string {
	// Para esta demostración, devolvemos un userID fijo
	// En una implementación real, decodificarías el JWT y extraerías el userID
	if len(token) > 10 {
		// Simular extracción de userID basada en el token
		return "user-001"
	}
	return ""
}