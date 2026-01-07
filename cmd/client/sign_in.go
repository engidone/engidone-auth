package main

import (
	"context"
	"log"

	pb "engidoneauth/internal/proto"
)

func signIn(client pb.AuthServiceClient, ctx context.Context, onError func(err error)) {
	// Llamar al servicio SignIn
	req := &pb.SignInRequest{
		Username: "caporras",
		Password: "password123",
	}

	resp, err := client.SignIn(ctx, req)
	if err != nil {
		onError(err)
		return
	}

	// Mostrar respuesta
	log.Printf("=== Respuesta del servicio SignIn ===")
	log.Printf("Token: %s", resp.Token)
	log.Printf("RefreshToken: %s", resp.RefreshToken)
	log.Printf("=====================================")
}
