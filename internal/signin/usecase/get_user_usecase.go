package usecase

import (
	"engidone-auth/internal/signin/domain"
)

// GetUserUseCase maneja la lógica de obtención de información de usuario
type GetUserUseCase struct {
	userRepo domain.UserRepository
}

// NewGetUserUseCase crea una nueva instancia del caso de uso de obtener usuario
func NewGetUserUseCase(userRepo domain.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
	}
}

// Execute ejecuta la obtención de información del usuario
func (uc *GetUserUseCase) Execute(userID string) (*domain.User, error) {
	// Validar userID
	if err := uc.validateUserID(userID); err != nil {
		return nil, err
	}

	// Obtener usuario del repositorio
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Preparar respuesta segura (sin información sensible)
	userResponse := &domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userResponse, nil
}

// ExecuteByUsername obtiene un usuario por su username
func (uc *GetUserUseCase) ExecuteByUsername(username string) (*domain.User, error) {
	// Validar username
	if err := uc.validateUsername(username); err != nil {
		return nil, err
	}

	// Obtener usuario del repositorio
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	// Preparar respuesta segura (sin información sensible)
	userResponse := &domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userResponse, nil
}

// validateUserID valida el userID de entrada
func (uc *GetUserUseCase) validateUserID(userID string) error {
	if userID == "" {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "El ID de usuario es requerido")
	}

	if len(userID) < 5 {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "ID de usuario inválido")
	}

	return nil
}

// validateUsername valida el username de entrada
func (uc *GetUserUseCase) validateUsername(username string) error {
	if username == "" {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "El nombre de usuario es requerido")
	}

	if len(username) < 3 {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "Nombre de usuario inválido")
	}

	return nil
}