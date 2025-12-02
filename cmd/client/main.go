package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "engidoneauth/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Simple flag to choose service
	if len(os.Args) < 2 {
		log.Fatal("Usage: client <signin|refresh>")
	}

	service := os.Args[1]

	// Connect to server
	conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	switch service {
	case "signin":
		signIn(client, ctx, func(err error) {})
	case "refresh":
		refreshToken(client, ctx, func(err error) {})
	case "hello":
		hello(client, ctx, func(err error) {})
	default:
		log.Fatalf("Invalid service: %s. Use 'signin' or 'refresh'", service)
	}
}
