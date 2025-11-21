package config

import (
	"engidoneauth/log"
)

func NewAppConfig(path string) *AppConfig {
	data, err := LoadFile[AppConfig](path + "/app.yaml")

	if err != nil {
		log.Fatal("Failed to load app config", log.Err(err))
		panic(err)
	}
	return data
}
