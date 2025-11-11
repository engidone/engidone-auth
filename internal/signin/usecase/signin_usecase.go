package usecase

import (
	"engidone-auth/internal/signin/domain"
)

// SigninUseCase maneja la lógica de autenticación de usuarios
type SigninUseCase struct {
	userRepo     domain.UserRepository
	tokenService domain.TokenService
}

// NewSigninUseCase crea una nueva instancia del caso de uso de signin
func NewSigninUseCase(userRepo domain.UserRepository, tokenService domain.TokenService) *SigninUseCase {
	return &SigninUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Execute ejecuta el proceso de autenticación
func (uc *SigninUseCase) Execute(credentials domain.Credentials) (*domain.AuthResponse, error) {
	// Validar credenciales
	if err := uc.validateCredentials(credentials); err != nil {
		return nil, err
	}

	// Verificar usuario y contraseña
	user, err := uc.userRepo.VerifyCredentials(credentials.Username, credentials.Password)
	if err != nil {
		return nil, err
	}

	// Generar token
	token, err := uc.tokenService.GenerateToken(user.ID)
	if err != nil {
		return nil, domain.NewAuthError(domain.ErrInvalidToken, "Error generando token de autenticación")
	}

	// Crear respuesta de autenticación
	response := &domain.AuthResponse{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Token:     token,
		ExpiresAt: domain.GetTokenExpiration(token),
	}

	return response, nil
}

// validateCredentials valida las credenciales de entrada
func (uc *SigninUseCase) validateCredentials(credentials domain.Credentials) error {
	if credentials.Username == "" {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "El nombre de usuario es requerido")
	}

	if credentials.Password == "" {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "La contraseña es requerida")
	}

	if len(credentials.Username) < 3 {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "El nombre de usuario debe tener al menos 3 caracteres")
	}

	if len(credentials.Password) < 4 {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "La contraseña debe tener al menos 4 caracteres")
	}

	return nil
}