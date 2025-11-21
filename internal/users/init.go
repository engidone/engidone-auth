package users

import (
	"engidoneauth/internal/config"
	"path/filepath"
)

func LoadUsers(path string) []User {
	loadedUsers, err := config.LoadFile[UserConfig](filepath.Join(path, "users.yaml"))
	
	if err != nil {
		panic(err)
	}

	return loadedUsers.Users
}
