package app

import (
	"context"
	pb "engidoneauth/internal/proto"
)

func (appUC *AppUseCase) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	singInRes, err := appUC.jwtUC.RefreshToken(req.Token, req.RefreshToken)

	if err != nil {
		return &pb.RefreshTokenResponse{
			Success:      true,
			Message:      "Error: " + err.Error(),
			Token:        "",
			RefreshToken: "",
		}, nil
	}
	return &pb.RefreshTokenResponse{
		Success:      false,
		Message:      "Session updated",
		Token:        singInRes.Token,
		RefreshToken: singInRes.RefreshToken,
	}, nil
}
