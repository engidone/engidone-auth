package app

import (
	"engidoneauth/internal/greet"
	"engidoneauth/internal/jwt"
	pb "engidoneauth/internal/proto"
	"engidoneauth/internal/signin"
)

type greetUseCase interface {
	SayHello(name string) (*greet.HelloResponse, error)
}

type signInUseCase interface {
	SingIn(credentials signin.Credentials) (*signin.Result, error)
}

type jwtUseCase interface {
	RefreshToken(token, refreshToken string) (*jwt.TokenInfo, error)
}

type AppUseCase struct {
	pb.UnimplementedAuthServiceServer
	greetUC  greetUseCase
	signInUC signInUseCase
	jwtUC    jwtUseCase
}

func NewUseCase(greetUC greetUseCase, signInUC signInUseCase, jwtUC jwtUseCase) pb.AuthServiceServer {
	return &AppUseCase{
		greetUC:  greetUC,
		signInUC: signInUC,
		jwtUC:    jwtUC,
	}
}
