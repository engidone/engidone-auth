package app

import (
	"engidoneauth/internal/greet"
	pb "engidoneauth/internal/proto"
	"engidoneauth/internal/signin"
)

type greetUseCase interface {
	SayHello(name string) (*greet.HelloResponse, error)
}

type signInUseCase interface {
	SingIn(credentials signin.Credentials) (*signin.Result, error)
}

type AppUseCase struct {
	pb.UnimplementedAuthServiceServer
	greetUC  greetUseCase
	signInUC signInUseCase
}

func NewUseCase(greetUC greetUseCase, signInUC signInUseCase) pb.AuthServiceServer {
	return &AppUseCase{
		greetUC:  greetUC,
		signInUC: signInUC,
	}
}
