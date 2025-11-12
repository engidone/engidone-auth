package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"engidone-auth/internal/database/seeders"

	_ "github.com/go-sql-driver/mysql"
)

// Migrator handles database migrations and seeders
type Migrator struct {
	db             *sql.DB
	migrationsPath string
	seedersPath    string
	logger         *log.Logger
	seederManager  *seeders.SeederManager
}

// NewMigrator creates a new migrator
func NewMigrator(db *sql.DB, migrationsPath string) *Migrator {
	logger := log.New(os.Stdout, "[MIGRATOR] ", log.LstdFlags)

	return &Migrator{
		db:             db,
		migrationsPath: migrationsPath,
		seedersPath:    filepath.Dir(migrationsPath) + "/seeders",
		logger:         logger,
	}
}

// NewMigratorWithSeeders creates a new migrator with custom seeders path
func NewMigratorWithSeeders(db *sql.DB, migrationsPath, seedersPath string) *Migrator {
	logger := log.New(os.Stdout, "[MIGRATOR] ", log.LstdFlags)

	return &Migrator{
		db:             db,
		migrationsPath: migrationsPath,
		seedersPath:    seedersPath,
		logger:         logger,
	}
}

// RunMigrations runs all pending migrations
func (m *Migrator) RunMigrations() error {
	m.logger.Println("Starting database migrations...")

	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("error creating migrations table: %w", err)
	}

	// Get migration files
	files, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("error getting migration files: %w", err)
	}

	// Get already executed migrations
	executedMigrations, err := m.getExecutedMigrations()
	if err != nil {
		return fmt.Errorf("error getting executed migrations: %w", err)
	}

	// Run pending migrations
	for _, file := range files {
		if !m.isMigrationExecuted(file, executedMigrations) {
			if err := m.runMigration(file); err != nil {
				return fmt.Errorf("error running migration %s: %w", file, err)
			}
		}
	}

	m.logger.Println("Migrations completed successfully")
	return nil
}

// RunMigrationsAndSeeders runs migrations and then seeders
func (m *Migrator) RunMigrationsAndSeeders() error {
	// Run migrations first
	if err := m.RunMigrations(); err != nil {
		return err
	}

	// Run seeders
	return m.RunSeeders()
}

// RunSeeders runs all configured seeders
func (m *Migrator) RunSeeders() error {
	m.logger.Println("Starting database seeders...")

	// Initialize seeder manager
	m.seederManager = seeders.NewSeederManager(m.db, m.seedersPath, m.logger)

	// Load seeders
	if err := m.seederManager.LoadSeeders(); err != nil {
		return fmt.Errorf("failed to load seeders: %w", err)
	}

	// Run seeders
	results, err := m.seederManager.RunSeeders()
	if err != nil {
		return fmt.Errorf("failed to run seeders: %w", err)
	}

	// Log results
	for _, result := range results {
		if result.Status == "success" {
			m.logger.Printf("✅ Seeder '%s' completed in %v", result.Name, result.Duration)
		} else {
			m.logger.Printf("❌ Seeder '%s' failed: %s", result.Name, result.Error)
		}
	}

	m.logger.Println("Seeders completed")
	return nil
}

// RunSeeder runs a specific seeder by name
func (m *Migrator) RunSeeder(name string) error {
	if m.seederManager == nil {
		m.seederManager = seeders.NewSeederManager(m.db, m.seedersPath, m.logger)
		if err := m.seederManager.LoadSeeders(); err != nil {
			return fmt.Errorf("failed to load seeders: %w", err)
		}
	}

	result := m.seederManager.RunSeeder(name)
	if result.Status == "success" {
		m.logger.Printf("✅ Seeder '%s' completed in %v", result.Name, result.Duration)
		return nil
	}

	return fmt.Errorf("seeder '%s' failed: %s", result.Name, result.Error)
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`

	_, err := m.db.Exec(query)
	return err
}

// getMigrationFiles gets all migration files sorted by filename
func (m *Migrator) getMigrationFiles() ([]string, error) {
	files, err := ioutil.ReadDir(m.migrationsPath)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

// getExecutedMigrations gets list of already executed migrations
func (m *Migrator) getExecutedMigrations() (map[string]bool, error) {
	rows, err := m.db.Query("SELECT filename FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		executed[filename] = true
	}

	return executed, rows.Err()
}

// isMigrationExecuted checks if a migration has been executed
func (m *Migrator) isMigrationExecuted(filename string, executed map[string]bool) bool {
	return executed[filename]
}

// runMigration executes a single migration file
func (m *Migrator) runMigration(filename string) error {
	// Read migration file
	filePath := filepath.Join(m.migrationsPath, filename)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration
	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query != "" && !strings.HasPrefix(query, "--") {
			if _, err := tx.Exec(query); err != nil {
				return fmt.Errorf("error executing query: %w", err)
			}
		}
	}

	// Record migration
	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (filename, executed_at) VALUES (?, ?)",
		filename, time.Now(),
	); err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}