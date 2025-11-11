package usecase

import (
	"engidone-auth/internal/signin/domain"
)

// RefreshTokenUseCase maneja la lógica de refresco de tokens
type RefreshTokenUseCase struct {
	userRepo     domain.UserRepository
	tokenService domain.TokenService
}

// NewRefreshTokenUseCase crea una nueva instancia del caso de uso de refresh token
func NewRefreshTokenUseCase(userRepo domain.UserRepository, tokenService domain.TokenService) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Execute ejecuta el proceso de refresco de token
func (uc *RefreshTokenUseCase) Execute(userID string, currentToken string) (*domain.AuthResponse, error) {
	// Validar userID
	if err := uc.validateUserID(userID); err != nil {
		return nil, err
	}

	// Validar token actual
	if err := uc.validateCurrentToken(currentToken); err != nil {
		return nil, err
	}

	// Verificar que el usuario existe
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, domain.NewAuthError(domain.ErrUserNotFound, "Usuario no encontrado")
	}

	// Refrescar el token
	tokenInfo, err := uc.tokenService.RefreshToken(currentToken)
	if err != nil {
		return nil, domain.NewAuthError(domain.ErrInvalidToken, "Error refrescando el token")
	}

	// Crear respuesta de autenticación actualizada
	response := &domain.AuthResponse{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Token:     tokenInfo.Token,
		ExpiresAt: tokenInfo.ExpiresAt,
	}

	return response, nil
}

// validateUserID valida el userID de entrada
func (uc *RefreshTokenUseCase) validateUserID(userID string) error {
	if userID == "" {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "El ID de usuario es requerido")
	}

	if len(userID) < 5 {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "ID de usuario inválido")
	}

	return nil
}

// validateCurrentToken valida el token actual
func (uc *RefreshTokenUseCase) validateCurrentToken(token string) error {
	if token == "" {
		return domain.NewAuthError(domain.ErrInvalidToken, "El token actual es requerido")
	}

	if len(token) < 10 {
		return domain.NewAuthError(domain.ErrInvalidToken, "Token actual inválido")
	}

	return nil
}