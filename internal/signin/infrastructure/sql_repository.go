package infrastructure

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"

	"engidone-auth/internal/database/sqlc"
	"engidone-auth/internal/signin/domain"
)

// SQLUserRepository implements UserRepository using SQL and SQLC
type SQLUserRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewSQLUserRepository creates a new SQL user repository
func NewSQLUserRepository(db *sql.DB) domain.UserRepository {
	return &SQLUserRepository{
		db:      db,
		queries: sqlc.New(),
	}
}

// FindByUsername busca un usuario por su username
func (r *SQLUserRepository) FindByUsername(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := r.queries.GetUserByUsername(ctx, r.db, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewAuthError(domain.ErrUserNotFound, "Usuario no encontrado")
		}
		return nil, domain.NewAuthError(domain.ErrInternalError, fmt.Sprintf("Error buscando usuario: %v", err))
	}

	return r.convertGetUserByUsernameToDomainUser(&user), nil
}

// FindByID busca un usuario por su ID
func (r *SQLUserRepository) FindByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert string ID to int32 for SQL
	var idInt int32
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		return nil, domain.NewAuthError(domain.ErrInvalidCredentials, "ID de usuario inválido")
	}

	user, err := r.queries.GetUserByID(ctx, r.db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewAuthError(domain.ErrUserNotFound, "Usuario no encontrado")
		}
		return nil, domain.NewAuthError(domain.ErrInternalError, fmt.Sprintf("Error buscando usuario: %v", err))
	}

	return r.convertGetUserByIDToDomainUser(&user), nil
}

// Create crea un nuevo usuario
func (r *SQLUserRepository) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Hash password
	hashedPassword := hashPassword(user.Password)

	params := sqlc.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}

	err := r.queries.CreateUser(ctx, r.db, params)
	if err != nil {
		// Check for duplicate entry errors
		if isDuplicateError(err) {
			if isDuplicateUsernameError(err) {
				return domain.NewAuthError("USER_EXISTS", "El nombre de usuario ya existe")
			}
			return domain.NewAuthError("EMAIL_EXISTS", "El email ya existe")
		}
		return domain.NewAuthError(domain.ErrInternalError, fmt.Sprintf("Error creando usuario: %v", err))
	}

	return nil
}

// Update actualiza un usuario existente
func (r *SQLUserRepository) Update(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert string ID to int32 for SQL
	var idInt int32
	if _, err := fmt.Sscanf(user.ID, "%d", &idInt); err != nil {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "ID de usuario inválido")
	}

	// Prepare update parameters
	params := sqlc.UpdateUserParams{
		ID:       idInt,
		Username: user.Username,
		Email:    user.Email,
	}

	// Only update password if provided
	if user.Password != "" {
		params.Password = hashPassword(user.Password)
	}

	err := r.queries.UpdateUser(ctx, r.db, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.NewAuthError(domain.ErrUserNotFound, "Usuario no encontrado")
		}
		if isDuplicateError(err) {
			return domain.NewAuthError("USER_EXISTS", "El nombre de usuario o email ya existe")
		}
		return domain.NewAuthError(domain.ErrInternalError, fmt.Sprintf("Error actualizando usuario: %v", err))
	}

	return nil
}

// Delete elimina un usuario por su ID
func (r *SQLUserRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert string ID to int32 for SQL
	var idInt int32
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		return domain.NewAuthError(domain.ErrInvalidCredentials, "ID de usuario inválido")
	}

	err := r.queries.DeleteUser(ctx, r.db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.NewAuthError(domain.ErrUserNotFound, "Usuario no encontrado")
		}
		return domain.NewAuthError(domain.ErrInternalError, fmt.Sprintf("Error eliminando usuario: %v", err))
	}

	return nil
}

// VerifyCredentials verifica las credenciales del usuario
func (r *SQLUserRepository) VerifyCredentials(username, password string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Hash the provided password
	hashedPassword := hashPassword(password)

	params := sqlc.VerifyCredentialsParams{
		Username: username,
		Password: hashedPassword,
	}

	user, err := r.queries.VerifyCredentials(ctx, r.db, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewAuthError(domain.ErrInvalidCredentials, "Credenciales inválidas")
		}
		return nil, domain.NewAuthError(domain.ErrInternalError, fmt.Sprintf("Error verificando credenciales: %v", err))
	}

	return r.convertVerifyCredentialsToDomainUser(&user), nil
}

// convertGetUserByUsernameToDomainUser converts GetUserByUsernameRow to domain User
func (r *SQLUserRepository) convertGetUserByUsernameToDomainUser(user *sqlc.GetUserByUsernameRow) *domain.User {
	domainUser := &domain.User{
		ID:       fmt.Sprintf("%d", user.ID),
		Username: user.Username,
		Email:    user.Email,
	}

	// Handle nullable timestamps
	if user.CreatedAt.Valid {
		domainUser.CreatedAt = user.CreatedAt.Time
	}
	if user.UpdatedAt.Valid {
		domainUser.UpdatedAt = user.UpdatedAt.Time
	}

	return domainUser
}

// convertGetUserByIDToDomainUser converts GetUserByIDRow to domain User
func (r *SQLUserRepository) convertGetUserByIDToDomainUser(user *sqlc.GetUserByIDRow) *domain.User {
	domainUser := &domain.User{
		ID:       fmt.Sprintf("%d", user.ID),
		Username: user.Username,
		Email:    user.Email,
	}

	// Handle nullable timestamps
	if user.CreatedAt.Valid {
		domainUser.CreatedAt = user.CreatedAt.Time
	}
	if user.UpdatedAt.Valid {
		domainUser.UpdatedAt = user.UpdatedAt.Time
	}

	return domainUser
}

// convertVerifyCredentialsToDomainUser converts VerifyCredentialsRow to domain User
func (r *SQLUserRepository) convertVerifyCredentialsToDomainUser(user *sqlc.VerifyCredentialsRow) *domain.User {
	domainUser := &domain.User{
		ID:       fmt.Sprintf("%d", user.ID),
		Username: user.Username,
		Email:    user.Email,
	}

	// Handle nullable timestamps
	if user.CreatedAt.Valid {
		domainUser.CreatedAt = user.CreatedAt.Time
	}
	if user.UpdatedAt.Valid {
		domainUser.UpdatedAt = user.UpdatedAt.Time
	}

	return domainUser
}

// isDuplicateError checks if the error is a duplicate entry error
func isDuplicateError(err error) bool {
	// MySQL duplicate entry error codes
	errMsg := err.Error()
	return contains(errMsg, "Duplicate entry") || contains(errMsg, "1062")
}

// isDuplicateUsernameError checks if the error is specifically for username duplicate
func isDuplicateUsernameError(err error) bool {
	errMsg := err.Error()
	return contains(errMsg, "username") || contains(errMsg, "users.username")
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
			 s[len(s)-len(substr):] == substr ||
			 indexOf(s, substr) >= 0)))
}

// indexOf finds the index of a substring
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// hashPassword generates a SHA-256 hash of the password
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}