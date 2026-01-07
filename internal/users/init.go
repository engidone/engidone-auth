package users

import (
	"path/filepath"

	"github.com/engidone/go-utils/common"
)

func LoadUsers(path string) []User {
	loadedUsers, err := common.LoadFile[UserConfig](filepath.Join(path, "users.yaml"))

	if err != nil {
		panic(err)
	}

	return loadedUsers.Users
}
