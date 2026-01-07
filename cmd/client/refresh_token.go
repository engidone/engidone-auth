package main

import (
	"context"
	"log"

	pb "engidoneauth/internal/proto"
)

var Token string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQwMjM0ODAsInN1YiI6IjRiNmNlNTM4LWMyNmItNDM5OS05OWRlLTI2MGE5ZTU5NGFiOCJ9.go7Tam8TdMkKgwtkcDrIL-yaQiMkbdp5zCr3bhertFRgAWtV6HafaJ-buFldY_kdbIB590bsa6kEWn00TlfCrl3-4rlA44CvgvZd5BkdhbgHzVcnCZ-P6w2OJaB_2M6_v2NKg7ziQlzhR53VytTFlbZrAESSsC3S-ko4AewQUq_up8B2rAT9SKnUY6qwdt8-nCC0fWrlIVIXEERGS73DGXFumPOAptX0ShLqbOX4sAAArZ7uYtHmAdQ-F3QUEoMqvWuqZHBKLdB6WYzWky1Ur6Q_re_vFcMei3jEwPi8Wj-F4PCHpBrqv7gsv6mgpLAKCyFQROj4RvRMUGzTWLEGwQ"
var RefreshToken string = "3Tx2i6ujIGHaEOaV-2dsZNvzxl7ntZGzgrs29Q9fpqLHrh8As_HESDUoPlNzAuUnZrz3CoMqzPnynfUQvRUcDA"

func refreshToken(client pb.AuthServiceClient, ctx context.Context, onError func(err error)) {
	// Llamar al servicio RefreshToken
	req := &pb.RefreshTokenRequest{
		RefreshToken: RefreshToken,
		Token:        Token,
	}

	resp, err := client.RefreshToken(ctx, req)
	if err != nil {
		onError(err)
		return
	}

	// Mostrar respuesta
	log.Printf("=== Respuesta del servicio RefreshToken ===")
	log.Printf("Token: %s", resp.Token)
	log.Printf("RefreshToken: %s", resp.RefreshToken)
	log.Printf("=====================================")
}
