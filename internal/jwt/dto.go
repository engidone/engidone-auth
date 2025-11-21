package jwt

import "crypto/rsa"

type TokenInfo struct {
	Token        string
	RefreshToken string
}

type Certs struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}
