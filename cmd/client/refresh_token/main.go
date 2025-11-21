package main

import (
	"context"
	"log"
	"time"

	pb "engidoneauth/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Token string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQwMjM0ODAsInN1YiI6IjRiNmNlNTM4LWMyNmItNDM5OS05OWRlLTI2MGE5ZTU5NGFiOCJ9.go7Tam8TdMkKgwtkcDrIL-yaQiMkbdp5zCr3bhertFRgAWtV6HafaJ-buFldY_kdbIB590bsa6kEWn00TlfCrl3-4rlA44CvgvZd5BkdhbgHzVcnCZ-P6w2OJaB_2M6_v2NKg7ziQlzhR53VytTFlbZrAESSsC3S-ko4AewQUq_up8B2rAT9SKnUY6qwdt8-nCC0fWrlIVIXEERGS73DGXFumPOAptX0ShLqbOX4sAAArZ7uYtHmAdQ-F3QUEoMqvWuqZHBKLdB6WYzWky1Ur6Q_re_vFcMei3jEwPi8Wj-F4PCHpBrqv7gsv6mgpLAKCyFQROj4RvRMUGzTWLEGwQ"
var RefreshToken string = "3Tx2i6ujIGHaEOaV-2dsZNvzxl7ntZGzgrs29Q9fpqLHrh8As_HESDUoPlNzAuUnZrz3CoMqzPnynfUQvRUcDA"

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
	req := &pb.RefreshTokenRequest{
		RefreshToken: RefreshToken,
		Token:        Token,
	}

	resp, err := client.RefreshToken(ctx, req)
	if err != nil {
		log.Fatalf("Error llamando al servicio RefreshToken: %v", err)
	}

	// Mostrar respuesta
	log.Printf("=== Respuesta del servicio RefreshToken ===")
	log.Printf("Mensaje: %s", resp.Message)
	log.Printf("Token: %s", resp.Token)
	log.Printf("RefreshToken: %s", resp.RefreshToken)
	log.Printf("=====================================")
}
