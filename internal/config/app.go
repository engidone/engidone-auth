package config

import (
	"github.com/engidone/go-utils/common"

	"github.com/engidone/go-utils/log"
)

func NewAppConfig(path string) *AppConfig {
	data, err := common.LoadFile[AppConfig](path + "/app.yaml")

	if err != nil {
		log.Fatal("Failed to load app config", log.Err(err))
		panic(err)
	}
	return data
}
