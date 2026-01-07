package signin

import "engidoneauth/internal/users"

type jwtUseCase interface {
	GenerateToken(userID string) (string, error)
	GetRefreshToken() (string, error)
	SyncRefreshToken(userID, refreshToken string) error
}

type credentialsUseCase interface {
	VerifyCredentials(userID string, password string) (bool, error)
}

type usersUseCase interface {
	GetUser(username string) (*users.User, error)
}

type UseCase struct {
	jwtUC         jwtUseCase
	credentialsUC credentialsUseCase
	usersUC       usersUseCase
}

func NewUseCase(jwtUC jwtUseCase, credentialsUC credentialsUseCase, usersUC usersUseCase) *UseCase {
	return &UseCase{
		jwtUC:         jwtUC,
		credentialsUC: credentialsUC,
		usersUC:       usersUC,
	}
}
