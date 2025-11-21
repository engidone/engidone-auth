package db

import (
	"database/sql"
	"engidoneauth/internal/config"
	"strings"

	_ "github.com/lib/pq"
)

func NewDBConnection(appConfig *config.AppConfig) (*sql.DB, error) {
	dsn := appConfig.Database.DSN

	dsn = strings.ReplaceAll(dsn, "{engine}", appConfig.Database.Engine)
	dsn = strings.ReplaceAll(dsn, "{user}", appConfig.Database.Username)
	dsn = strings.ReplaceAll(dsn, "{password}", appConfig.Database.Password)
	dsn = strings.ReplaceAll(dsn, "{host}", appConfig.Database.Host)
	dsn = strings.ReplaceAll(dsn, "{port}", appConfig.Database.Port)
	dsn = strings.ReplaceAll(dsn, "{db_name}", appConfig.Database.DBName)
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
