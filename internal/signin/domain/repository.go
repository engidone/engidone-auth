package domain

// UserRepository define la interfaz para el repositorio de usuarios
type UserRepository interface {
	// FindByUsername busca un usuario por su username
	FindByUsername(username string) (*User, error)

	// FindByID busca un usuario por su ID
	FindByID(id string) (*User, error)

	// Create crea un nuevo usuario
	Create(user *User) error

	// Update actualiza un usuario existente
	Update(user *User) error

	// Delete elimina un usuario por su ID
	Delete(id string) error

	// VerifyCredentials verifica las credenciales del usuario
	VerifyCredentials(username, password string) (*User, error)
}

// Use case interfaces for GoKit
type SigninUseCase interface {
	Execute(credentials Credentials) (*AuthResponse, error)
}

type ValidateTokenUseCase interface {
	Execute(token string) (*User, error)
}

type RefreshTokenUseCase interface {
	Execute(userID, token string) (*AuthResponse, error)
}

type GetUserUseCase interface {
	Execute(userID string) (*User, error)
}