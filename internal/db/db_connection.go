package db

import (
	"database/sql"
	"engidoneauth/internal/config"
	"engidoneauth/util/env"
	"strings"

	_ "github.com/lib/pq"
)

func NewDBConnection(appConfig *config.AppConfig) (*sql.DB, error) {
	dsn := appConfig.Database.DSN

	dsn = strings.ReplaceAll(dsn, "{engine}", appConfig.Database.Engine)
	dsn = strings.ReplaceAll(dsn, "{user}", env.Get("DB_USER"))
	dsn = strings.ReplaceAll(dsn, "{password}", env.Get("DB_PASSWORD"))
	dsn = strings.ReplaceAll(dsn, "{host}", env.Get("DB_HOST"))
	dsn = strings.ReplaceAll(dsn, "{port}", env.Get("DB_PORT"))
	dsn = strings.ReplaceAll(dsn, "{db_name}", env.Get("DB_NAME"))
	dsn = strings.ReplaceAll(dsn, "{ssl_mode}", appConfig.Database.SSLMode)
	db, err := sql.Open(appConfig.Database.Engine, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
