package server

import (
	"context"
	"database/sql"

	"go.uber.org/fx"
)

func registerDBConnectionLoader(lc fx.Lifecycle, dbConn *sql.DB) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// El ping valida la conexión
			return dbConn.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return dbConn.Close()
		},
	})
}
