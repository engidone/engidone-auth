package users

type repository interface {
	findUserByUsername(username string) (*User, error)
}

type UseCase struct {
	repo repository
}

func NewUseCase(userRepository repository) *UseCase {
	return &UseCase{
		repo: userRepository,
	}
}
