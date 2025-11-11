package di

import (
	"go.uber.org/fx"

	"engidone-auth/internal/signin/domain"
	"engidone-auth/internal/signin/infrastructure"
	"engidone-auth/internal/signin/usecase"
)

// SigninModule provides all signin service dependencies
var SigninModule = fx.Options(
	fx.Provide(
		NewUserRepository,
		NewTokenService,
		NewSigninUseCase,
		NewValidateTokenUseCase,
		NewRefreshTokenUseCase,
		NewGetUserUseCase,
	),
)

// NewUserRepository provides a UserRepository implementation
func NewUserRepository() domain.UserRepository {
	return infrastructure.NewMemoryUserRepository()
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