package app

import (
	"context"
	pb "engidoneauth/internal/proto"
)

func (uc *AppUseCase) Hello(cxt context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Name:    request.Name,
		Message: "Hola " + request.Name + " Desde el servidor GRPC Docker v6",
	}, nil
}
