package transport

import (
	"context"

	"engidone-auth/internal/hello/endpoints"
	pb "engidone-auth/internal/hello/proto"
)

type grpcServer struct {
	pb.UnimplementedHelloServiceServer
	endpoints endpoints.Set
}

func NewGRPCServer(endpoints endpoints.Set) pb.HelloServiceServer {
	return &grpcServer{
		endpoints: endpoints,
	}
}

func (g *grpcServer) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	request := endpoints.HelloRequest{
		Name: req.Name,
	}

	response, err := g.endpoints.HelloEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	resp := response.(endpoints.HelloResponse)
	return &pb.HelloResponse{
		Message: resp.Message,
		Success: resp.Success,
	}, nil
}