package jwt

import "time"

type UseCase struct {
	tokenDuration time.Duration
}

func NewUseCase(tokenDuration time.Duration) *UseCase {
	return &UseCase{tokenDuration: tokenDuration}
}
