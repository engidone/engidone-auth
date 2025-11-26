package recovery

type repository interface {
	findRecoveryCode(code string) (string, error)
}

type UseCase struct {
	repo repository
}

func NewUseCase(recoveryRepository repository) *UseCase {
	return &UseCase{
		repo: recoveryRepository,
	}
}
