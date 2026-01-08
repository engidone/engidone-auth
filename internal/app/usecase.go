package app

import (
	"engidoneauth/internal/jwt"
	pb "engidoneauth/internal/proto"
	"engidoneauth/internal/signin"
)

type signInUseCase interface {
	SingIn(credentials signin.Credentials) (*signin.Result, error)
}

type jwtUseCase interface {
	RefreshToken(token, refreshToken string) (*jwt.TokenInfo, error)
}

type AppUseCase struct {
	pb.UnimplementedAuthServiceServer
	signInUC signInUseCase
	jwtUC    jwtUseCase
}

func NewUseCase(signInUC signInUseCase, jwtUC jwtUseCase) pb.AuthServiceServer {
	return &AppUseCase{
		signInUC: signInUC,
		jwtUC:    jwtUC,
	}
}
