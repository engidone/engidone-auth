package app

import (
	"context"

	"engidoneauth/internal/jwt"
	pb "engidoneauth/internal/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (appUC *AppUseCase) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	singInRes, err := appUC.jwtUC.RefreshToken(req.Token, req.RefreshToken)

	if err != nil {
		// Simple error mapping for refresh token
		if jwt.IsErrorCode(err, jwt.CodeInvalidToken) || jwt.IsErrorCode(err, jwt.CodeInvalidRefreshToken) {
			return nil, status.Error(codes.Unauthenticated, "Invalid or expired session")
		}

		if jwt.IsErrorCode(err, jwt.CodeTokenExpired) {
			return nil, status.Error(codes.Unauthenticated, "Session expired, please login again")
		}

		if jwt.IsErrorCode(err, jwt.CodeTokenGeneration) || jwt.IsErrorCode(err, jwt.CodeGeneratingRefreshToken) {
			return nil, status.Error(codes.Internal, "Token service temporarily unavailable")
		}

		// Generic error
		return nil, status.Error(codes.Internal, "Session refresh failed")
	}

	return &pb.RefreshTokenResponse{
		Token:        singInRes.Token,
		RefreshToken: singInRes.RefreshToken,
	}, nil
}
