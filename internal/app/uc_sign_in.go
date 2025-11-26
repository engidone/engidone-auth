package app

import (
	"context"

	"engidoneauth/internal/credentials"
	pb "engidoneauth/internal/proto"
	"engidoneauth/internal/signin"
	"engidoneauth/internal/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (appUC *AppUseCase) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	singInRes, err := appUC.signInUC.SingIn(signin.Credentials{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		// Simple error mapping to client-friendly responses
		if credentials.IsErrorCode(err, credentials.CodeInvalidCredentials) ||
			users.IsErrorCode(err, users.CodeUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, "Invalid username or password")
		}

		if signin.IsErrorCode(err, signin.CodeMissingUsername) {
			return nil, status.Error(codes.InvalidArgument, "Username is required")
		}

		if signin.IsErrorCode(err, signin.CodeMissingPassword) {
			return nil, status.Error(codes.InvalidArgument, "Password is required")
		}

		if signin.IsErrorCode(err, signin.CodeUsernameTooShort) {
			return nil, status.Error(codes.InvalidArgument, "Username must be at least 3 characters long")
		}

		if signin.IsErrorCode(err, signin.CodePasswordTooShort) {
			return nil, status.Error(codes.InvalidArgument, "Password must be at least 4 characters long")
		}

		// Generic error for unknown issues
		return nil, status.Error(codes.Internal, "Authentication service temporarily unavailable")
	}

	return &pb.SignInResponse{
		Token:        singInRes.Token,
		RefreshToken: singInRes.RefreshToken,
	}, nil
}
