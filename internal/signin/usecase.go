package signin

import "engidoneauth/internal/users"

type jwtUseCase interface {
	GenerateToken(userID int32) (string, error)
}

type credentialsUseCase interface {
	VerifyCredentials(userID int32, password string) (bool, error)
}

type usersUsecase interface {
	GetUser(username string) (*users.User, error)
}

type repository interface {
}

type UseCase struct {
	repository    repository
	jwtUC         jwtUseCase
	credentialsUC credentialsUseCase
	usersUC       usersUsecase
}

func NewUseCase(repo repository, jwtUC jwtUseCase, credentialsUC credentialsUseCase, usersUC usersUsecase) *UseCase {
	return &UseCase{
		repository:    repo,
		jwtUC:         jwtUC,
		credentialsUC: credentialsUC,
		usersUC:       usersUC,
	}
}
