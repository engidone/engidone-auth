package usecase

import (
	"engidone-auth/internal/signin/domain"
)

// ValidateTokenUseCase maneja la lógica de validación de tokens
type ValidateTokenUseCase struct {
	userRepo     domain.UserRepository
	tokenService domain.TokenService
}

// NewValidateTokenUseCase crea una nueva instancia del caso de uso de validación de token
func NewValidateTokenUseCase(userRepo domain.UserRepository, tokenService domain.TokenService) *ValidateTokenUseCase {
	return &ValidateTokenUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Execute ejecuta la validación del token
func (uc *ValidateTokenUseCase) Execute(token string) (*domain.User, error) {
	// Validar formato del token
	if err := uc.validateTokenFormat(token); err != nil {
		return nil, err
	}

	// Validar y extraer información del token
	tokenInfo, err := uc.tokenService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Verificar que el usuario existe
	user, err := uc.userRepo.FindByID(tokenInfo.UserID)
	if err != nil {
		return nil, domain.NewAuthError(domain.ErrUserNotFound, "Usuario del token no encontrado")
	}

	return user, nil
}

// validateTokenFormat valida el formato básico del token
func (uc *ValidateTokenUseCase) validateTokenFormat(token string) error {
	if token == "" {
		return domain.NewAuthError(domain.ErrInvalidToken, "El token es requerido")
	}

	if len(token) < 10 {
		return domain.NewAuthError(domain.ErrInvalidToken, "Formato de token inválido")
	}

	// Validar que comience con "Bearer " (formato estándar)
	if len(token) >= 7 && token[:7] == "Bearer " {
		return nil
	}

	return domain.NewAuthError(domain.ErrInvalidToken, "Formato de token inválido. Debe comenzar con 'Bearer '")
}