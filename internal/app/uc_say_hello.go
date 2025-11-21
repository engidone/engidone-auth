package app

import (
	"context"
	pb "engidoneauth/internal/proto"
)

func (appUC *AppUseCase) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	res, err := appUC.greetUC.SayHello(req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.HelloResponse{
		Message: res.Message,
	}, nil
}
