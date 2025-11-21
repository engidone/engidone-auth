package app

import (
	"context"
	pb "engidoneauth/internal/proto"
	"engidoneauth/internal/signin"
)

func (appUC *AppUseCase) RefreshToken(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	singInRes, err := appUC.signInUC.SingIn(signin.Credentials{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return &pb.SignInResponse{
			Message:      "Opps " + err.Error(),
			Token:        "Yuca",
			RefreshToken: "Yuca",
		}, nil
	}
	return &pb.SignInResponse{
		Message:      "Welcome " + req.Username,
		Token:        singInRes.Token,
		RefreshToken: singInRes.RefreshToken,
	}, nil
}
