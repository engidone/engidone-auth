package transport

import (
	"context"

	"engidone-auth/internal/signin/endpoints"
	pb "engidone-auth/internal/signin/proto"
)

type grpcServer struct {
	pb.UnimplementedSigninServiceServer
	endpoints endpoints.Set
}

func NewGRPCServer(endpoints endpoints.Set) pb.SigninServiceServer {
	return &grpcServer{
		endpoints: endpoints,
	}
}

func (g *grpcServer) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	request := endpoints.SigninRequest{
		Username: req.Username,
		Password: req.Password,
	}

	response, err := g.endpoints.SigninEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	resp := response.(endpoints.SigninResponse)
	return &pb.SigninResponse{
		Success:   resp.Success,
		Message:   resp.Message,
		UserId:    resp.UserID,
		Username:  resp.Username,
		Email:     resp.Email,
		Token:     resp.Token,
		ExpiresAt: resp.ExpiresAt,
	}, nil
}

func (g *grpcServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	request := endpoints.ValidateTokenRequest{
		Token: req.Token,
	}

	response, err := g.endpoints.ValidateTokenEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	resp := response.(endpoints.ValidateTokenResponse)
	return &pb.ValidateTokenResponse{
		Valid:    resp.Valid,
		Message:  resp.Message,
		UserId:   resp.UserID,
		Username: resp.Username,
		Email:    resp.Email,
	}, nil
}

func (g *grpcServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.SigninResponse, error) {
	request := endpoints.RefreshTokenRequest{
		UserID: req.UserId,
		Token:  req.Token,
	}

	response, err := g.endpoints.RefreshTokenEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	resp := response.(endpoints.SigninResponse)
	return &pb.SigninResponse{
		Success:   resp.Success,
		Message:   resp.Message,
		UserId:    resp.UserID,
		Username:  resp.Username,
		Email:     resp.Email,
		Token:     resp.Token,
		ExpiresAt: resp.ExpiresAt,
	}, nil
}

func (g *grpcServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	request := endpoints.GetUserRequest{
		UserID: req.UserId,
	}

	response, err := g.endpoints.GetUserEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	resp := response.(endpoints.GetUserResponse)
	return &pb.GetUserResponse{
		Success:   resp.Success,
		Message:   resp.Message,
		UserId:    resp.UserID,
		Username:  resp.Username,
		Email:     resp.Email,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}, nil
}