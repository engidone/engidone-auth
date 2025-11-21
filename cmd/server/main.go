package main

import (
	c "engidoneauth/common"
	"engidoneauth/internal/config"
	"engidoneauth/internal/jwt"
	"engidoneauth/internal/server"
	"engidoneauth/internal/users"

	"engidoneauth/log"
	"path/filepath"
)

func main() {
	paths := config.NewConfigPaths(c.BACK, c.BACK)
	cf := config.NewAppConfig(paths.Config)

	publicKey, err := jwt.LoadPublicKey(
		filepath.Join(paths.Root, cf.Certs.Public),
	)

	if err != nil {
		log.Fatalf(jwt.ErrLoadingPublicKey, err.Error())
	}

	privateKey, err := jwt.LoadPrivateKey(
		filepath.Join(paths.Root, cf.Certs.Private),
	)

	if err != nil {
		log.Fatalf(jwt.ErrLoadingPublicKey, err.Error())
	}

	certs := jwt.Certs{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	userList := users.LoadUsers(paths.Config)

	server.NewGRPCServer(cf, certs, userList)
}
