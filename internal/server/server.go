package server

import (
	"context"
	"engidoneauth/internal/config"
	pb "engidoneauth/internal/proto"
	"fmt"
	"log"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func newGRPCServer(lc fx.Lifecycle, authSvr pb.AuthServiceServer, appConfig *config.AppConfig) *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authSvr)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfig.Application.Port))
			if err != nil {
				return err
			}
			go func() {
				log.Printf("gRPC server listening on :%s", appConfig.Application.Port)
				if err := grpcServer.Serve(lis); err != nil {
					log.Fatalf("Failed to serve: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})

	return grpcServer
}