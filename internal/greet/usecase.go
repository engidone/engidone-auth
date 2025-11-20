package greet

type repository interface {
	sendGreeting(name string) (string, error)
}
// HelloUC implements the greeting use case
type UseCase struct {
	repository repository
}

// NewHelloUseCase creates a new hello use case
func NewUseCase(repository repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

