package main

import (
	"context"
	"log"
	"time"

	pb "engidoneauth/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Conectar al servidor gRPC
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Llamar al servicio Hello
	req := &pb.SignInRequest{
		Username: "caporras",
		Password: "password123",
	}

	resp, err := client.SignIn(ctx, req)
	if err != nil {
		log.Fatalf("Error llamando al servicio SignIn: %v", err)
	}

	// Mostrar respuesta
	log.Printf("=== Respuesta del servicio SignIn ===")
	log.Printf("Mensaje: %s", resp.Message)
	log.Printf("Token: %s", resp.Token)
	log.Printf("RefreshToken: %s", resp.RefreshToken)
	log.Printf("=====================================")
}
