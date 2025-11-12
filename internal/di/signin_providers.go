package di

import (
	"database/sql"

	"go.uber.org/fx"

	"engidone-auth/internal/database"
	"engidone-auth/internal/signin/domain"
	"engidone-auth/internal/signin/infrastructure"
	"engidone-auth/internal/signin/usecase"
)

// SigninModule provides all signin service dependencies
var SigninModule = fx.Options(
	fx.Provide(
		database.NewConfig,
		database.NewConnection,
		NewUserRepository,
		NewTokenService,
		NewSigninUseCase,
		NewValidateTokenUseCase,
		NewRefreshTokenUseCase,
		NewGetUserUseCase,
	),
	fx.Invoke(RunMigrationsAndSeeders),
)

// NewUserRepository provides a UserRepository implementation using SQL
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return infrastructure.NewSQLUserRepository(db)
}

// RunMigrations runs database migrations on application startup
func RunMigrations(db *sql.DB) error {
	migrator := database.NewMigrator(db, "internal/database/migrations")
	return migrator.RunMigrations()
}

// RunMigrationsAndSeeders runs migrations and seeders on application startup
func RunMigrationsAndSeeders(db *sql.DB) error {
	migrator := database.NewMigrator(db, "internal/database/migrations")
	return migrator.RunMigrationsAndSeeders()
}

// NewTokenService provides a TokenService implementation
func NewTokenService() domain.TokenService {
	return domain.NewJWTTokenService()
}

// NewSigninUseCase provides a SigninUseCase implementation
func NewSigninUseCase(userRepo domain.UserRepository, tokenService domain.TokenService) domain.SigninUseCase {
	return usecase.NewSigninUseCase(userRepo, tokenService)
}

// NewValidateTokenUseCase provides a ValidateTokenUseCase implementation
func NewValidateTokenUseCase(userRepo domain.UserRepository, tokenService domain.TokenService) domain.ValidateTokenUseCase {
	return usecase.NewValidateTokenUseCase(userRepo, tokenService)
}

// NewRefreshTokenUseCase provides a RefreshTokenUseCase implementation
func NewRefreshTokenUseCase(userRepo domain.UserRepository, tokenService domain.TokenService) domain.RefreshTokenUseCase {
	return usecase.NewRefreshTokenUseCase(userRepo, tokenService)
}

// NewGetUserUseCase provides a GetUserUseCase implementation
func NewGetUserUseCase(userRepo domain.UserRepository) domain.GetUserUseCase {
	return usecase.NewGetUserUseCase(userRepo)
}