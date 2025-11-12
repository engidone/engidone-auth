package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

// NewConfig creates database configuration from environment variables
func NewConfig() *Config {
	config := &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		Database: getEnv("DB_NAME", "engidone_auth"),
		Username: getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
	}

	return config
}

// NewConnection creates a new database connection
func NewConnection(config *Config) (*sql.DB, error) {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&collation=utf8mb4_unicode_ci",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}