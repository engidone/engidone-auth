package config

import (
	"engidoneauth/log"
)

func newAppConfig(paths Paths) *AppConfig {
	data, err := LoadFile[AppConfig](paths.Config + "/app.yaml")
	if err != nil {
		log.Fatal("Failed to load app config", log.Err(err))
		panic(err)
	}
	return data
}
