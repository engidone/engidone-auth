package db

import (
	"database/sql"
	"engidoneauth/internal/config"

	"github.com/engidone/go-utils/database"

	_ "github.com/lib/pq"
)

func NewDBConnection(appConfig *config.AppConfig) (*sql.DB, error) {
	dsn := appConfig.Database.DSN

	dsn = database.BuildDSN(appConfig.Database)

	db, err := sql.Open(appConfig.Database.Engine, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
