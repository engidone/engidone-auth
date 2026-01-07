package main

import (
	"engidoneauth/internal/config"
	"engidoneauth/internal/jwt"
	"engidoneauth/internal/server"
	"engidoneauth/internal/users"

	"github.com/engidone/go-utils/common"

	"path/filepath"

	"github.com/engidone/go-utils/log"
)

func main() {
	paths := common.NewConfigPaths("cmd/config")
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
