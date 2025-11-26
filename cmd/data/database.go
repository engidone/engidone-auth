package main

import (
	"database/sql"
	"engidoneauth/internal/config"
	"engidoneauth/log"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/lib/pq"
)

// openDb creates and returns a database connection using the app configuration
func openDb() (*sql.DB, error) {
	appConfig, err := config.LoadFile[config.AppConfig](getConfigPath("config", "app.yaml"))
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	database := appConfig.Database
	dsn := buildDSN(database)

	db, err := sql.Open(database.Engine, dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	return db, nil
}

// DatabaseConfig represents the database configuration structure
type DatabaseConfig struct {
	DSN      string
	Engine   string
	Username string
	Password string
	Host     string
	Port     string
	SSLMode  string
	DBName   string
}

// buildDSN constructs the database connection string from config
func buildDSN(database struct {
	DSN      string `yaml:"dsn"`
	Engine   string `yaml:"engine"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	SSLMode  string `yaml:"ssl_mode"`
	DBName   string `yaml:"db_name"`
}) string {
	dsn := database.DSN
	dsn = strings.ReplaceAll(dsn, "{engine}", database.Engine)
	dsn = strings.ReplaceAll(dsn, "{user}", database.Username)
	dsn = strings.ReplaceAll(dsn, "{password}", database.Password)
	dsn = strings.ReplaceAll(dsn, "{host}", database.Host)
	dsn = strings.ReplaceAll(dsn, "{port}", database.Port)
	dsn = strings.ReplaceAll(dsn, "{db_name}", database.DBName)
	dsn = strings.ReplaceAll(dsn, "{ssl_mode}", database.SSLMode)
	return dsn
}

// truncateTables truncates the specified database tables
func truncateTables(db *sql.DB, tableNames []string) {
	fmt.Println("🧹 Starting table truncation...")

	for _, tableName := range tableNames {
		truncateSingleTable(db, tableName)
	}

	restoreForeignKeyConstraints(db)
	fmt.Println("🎉 Table truncation completed")
}

// truncateSingleTable truncates a single table and shows count
func truncateSingleTable(db *sql.DB, tableName string) {
	// Disable foreign key constraints temporarily
	disableFKQuery := "SET session_replication_role = replica;"
	_, err := db.Exec(disableFKQuery)
	if err != nil {
		log.Fatalf("Error desactivando restricciones de clave externa: %v", err)
	}

	// Count records before truncating
	showTableRecordCount(db, tableName)

	// Truncate table
	truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName)
	_, err = db.Exec(truncateQuery)
	if err != nil {
		log.Fatalf("Error truncando tabla %s: %v", tableName, err)
	}

	fmt.Printf("✅ Table %s truncated successfully\n", tableName)
}

// showTableRecordCount displays the number of records in a table
func showTableRecordCount(db *sql.DB, tableName string) {
	var count int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.QueryRow(countQuery).Scan(&count)
	if err != nil {
		fmt.Printf("⚠️  Could not count records in table %s: %v\n", tableName, err)
	} else {
		fmt.Printf("📊 Table %s: %d records found\n", tableName, count)
	}
}

// restoreForeignKeyConstraints restores normal foreign key constraint checking
func restoreForeignKeyConstraints(db *sql.DB) {
	enableFKQuery := "SET session_replication_role = DEFAULT;"
	_, err := db.Exec(enableFKQuery)
	if err != nil {
		log.Fatalf("Error restaurando restricciones de clave externa: %v", err)
	}
}

// getConfigPath returns the absolute path to a config file relative to the caller
func getConfigPath(args ...string) string {
	_, b, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(b), "..", filepath.Join(args...))
}

// getDirPath returns the directory path of the calling file
func getDirPath(b string) string {
	return filepath.Dir(b)
}

// getCallerDirPath returns the directory path of the current caller
func getCallerDirPath() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("Could not get file path")
	}
	return filepath.Dir(filename)
}