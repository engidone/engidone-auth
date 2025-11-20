package users

import (
	"context"
	"engidoneauth/internal/config"

	"go.uber.org/fx"
)

func registerUsersLoader(lc fx.Lifecycle, users *[]User, paths config.Paths) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			loadedUsers, err := config.LoadFile[UserConfig](paths.Config + "/users.yaml")
			if err != nil {
				return err
			}
			*users = loadedUsers.Users
			return nil
		},
	})
}
