package infrastructure

import (
	"crypto/sha256"
	"fmt"
	"time"

	"engidone-auth/internal/signin/domain"
)

// MemoryUserRepository implementa UserRepository en memoria
type MemoryUserRepository struct {
	users map[string]*domain.User
}

// NewMemoryUserRepository crea una nueva instancia del repositorio en memoria
func NewMemoryUserRepository() *MemoryUserRepository {
	repo := &MemoryUserRepository{
		users: make(map[string]*domain.User),
	}

	// Inicializar con usuarios quemados (hash de contraseñas)
	repo.seedUsers()
	return repo
}

// seedUsers inserta usuarios quemados para demostración
func (r *MemoryUserRepository) seedUsers() {
	now := time.Now()

	// Usuario admin: password123
	adminUser := &domain.User{
		ID:        "user-001",
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  hashPassword("password123"),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Usuario test: test123
	testUser := &domain.User{
		ID:        "user-002",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashPassword("test123"),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Usuario john: john123
	johnUser := &domain.User{
		ID:        "user-003",
		Username:  "john",
		Email:     "john@example.com",
		Password:  hashPassword("john123"),
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.users["admin"] = adminUser
	r.users["testuser"] = testUser
	r.users["john"] = johnUser
}

// hashPassword genera un hash SHA-256 de la contraseña
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// FindByUsername busca un usuario por su username
func (r *MemoryUserRepository) FindByUsername(username string) (*domain.User, error) {
	user, exists := r.users[username]
	if !exists {
		return nil, domain.NewAuthError(domain.ErrUserNotFound, "Usuario no encontrado")
	}

	// Devolver una copia sin la contraseña
	userCopy := &domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userCopy, nil
}

// FindByID busca un usuario por su ID
func (r *MemoryUserRepository) FindByID(id string) (*domain.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			// Devolver una copia sin la contraseña
			userCopy := &domain.User{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}
			return userCopy, nil
		}
	}

	return nil, domain.NewAuthError(domain.ErrUserNotFound, "Usuario no encontrado")
}

// Create crea un nuevo usuario
func (r *MemoryUserRepository) Create(user *domain.User) error {
	if _, exists := r.users[user.Username]; exists {
		return domain.NewAuthError("USER_EXISTS", "El usuario ya existe")
	}

	user.Password = hashPassword(user.Password)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.users[user.Username] = user
	return nil
}

// Update actualiza un usuario existente
func (r *MemoryUserRepository) Update(user *domain.User) error {
	if _, exists := r.users[user.Username]; !exists {
		return domain.NewAuthError(domain.ErrUserNotFound, "Usuario no encontrado")
	}

	existingUser := r.users[user.Username]
	existingUser.Email = user.Email
	if user.Password != "" {
		existingUser.Password = hashPassword(user.Password)
	}
	existingUser.UpdatedAt = time.Now()

	return nil
}

// Delete elimina un usuario por su ID
func (r *MemoryUserRepository) Delete(id string) error {
	for username, user := range r.users {
		if user.ID == id {
			delete(r.users, username)
			return nil
		}
	}

	return domain.NewAuthError(domain.ErrUserNotFound, "Usuario no encontrado")
}

// VerifyCredentials verifica las credenciales del usuario
func (r *MemoryUserRepository) VerifyCredentials(username, password string) (*domain.User, error) {
	user, exists := r.users[username]
	if !exists {
		return nil, domain.NewAuthError(domain.ErrUserNotFound, "Usuario no encontrado")
	}

	// Verificar contraseña
	hashedPassword := hashPassword(password)
	if user.Password != hashedPassword {
		return nil, domain.NewAuthError(domain.ErrInvalidCredentials, "Credenciales inválidas")
	}

	// Devolver una copia sin la contraseña
	userCopy := &domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userCopy, nil
}