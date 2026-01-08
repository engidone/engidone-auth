package main

import (
	"context"
	"log"

	pb "engidoneauth/internal/proto"
)

func hello(client pb.AuthServiceClient, ctx context.Context, onError func(err error)) {
	req := &pb.HelloRequest{
		Name: "Carlos",
	}

	resp, err := client.Hello(ctx, req)
	if err != nil {
		onError(err)
		return
	}

	log.Printf("=== Respuesta del servicio Hello ===")
	log.Printf("Nombre: %s", resp.Name)
	log.Printf("Mensaje: %s", resp.Message)
	log.Printf("=====================================")
}
