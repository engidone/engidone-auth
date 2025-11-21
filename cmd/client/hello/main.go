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
	// Conectar al servidor gRPC
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	// Solicitar nombre si no se proporciona como argumento
	name := "Carlos wey"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Llamar al servicio Hello
	req := &pb.HelloRequest{
		Name: name,
	}

	resp, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("Error llamando al servicio Hello: %v", err)
	}

	// Mostrar respuesta
	log.Printf("=== Respuesta del servicio Hello ===")
	log.Printf("Nombre: %s", name)
	log.Printf("Mensaje: %s", resp.Message)
	log.Printf("=====================================")
}
