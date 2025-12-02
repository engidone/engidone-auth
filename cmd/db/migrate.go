package main

import (
	"database/sql"
	"engidoneauth/internal/config"
	"engidoneauth/util/env"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"engidoneauth/log"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

var rollback bool
var steps int64
var sequential int64

func main() {
	var rootCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Execute database migrations",
		Long: `Execute database migrations with various options:

Examples:
  migrate                    # Apply all pending migrations
  migrate -r                 # Rollback one migration
  migrate -r -k 3            # Rollback 3 migrations
  migrate -s 5               # Apply up to migration 00005
  migrate -r -s 3            # Rollback to migration 00003`,
		Run: func(cmd *cobra.Command, args []string) {
			if sequential > 0 {
				if rollback {
					runRollbackToSequential(sequential)
				} else {
					runMigrateToSequential(sequential)
				}
			} else if rollback {
				if steps > 0 {
					runRollback(steps)
				} else {
					runRollback(1) // default rollback 1 step
				}
			} else {
				runMigration()
			}
		},
	}

	rootCmd.Flags().BoolVarP(&rollback, "rollback", "r", false, "Rollback migrations")
	rootCmd.Flags().Int64VarP(&steps, "steps", "k", 1, "Number of migration steps to rollback")
	rootCmd.Flags().Int64VarP(&sequential, "sequential", "s", 0, "Specific sequential migration number (e.g., 5 for 00005)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

func getConfigPath(args ...string) string {
	_, b, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(b), "..", filepath.Join(args...))
}

func resolveMigrationsPath() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b) // directorio cmd/
	// Ajusta ruta de migrations según estructura del proyecto desde cmd/
	return filepath.Join(basepath, "../../migrations")
}

func openDb() (*sql.DB, error) {
	appConfig, err := config.LoadFile[config.AppConfig](getConfigPath("config", "app.yaml"))

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	database := appConfig.Database
	dsn := database.DSN
	dsn = strings.ReplaceAll(dsn, "{engine}", database.Engine)
	dsn = strings.ReplaceAll(dsn, "{user}", env.Get("DB_USER"))
	dsn = strings.ReplaceAll(dsn, "{password}", env.Get("DB_PASSWORD"))
	dsn = strings.ReplaceAll(dsn, "{host}", env.Get("DB_HOST"))
	dsn = strings.ReplaceAll(dsn, "{port}", env.Get("DB_PASSWORD"))
	dsn = strings.ReplaceAll(dsn, "{db_name}", env.Get("DB_NAME"))
	dsn = strings.ReplaceAll(dsn, "{ssl_mode}", database.SSLMode)

	db, err := sql.Open(database.Engine, dsn)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	return db, nil
}

func runMigration() {
	db, err := openDb()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	migrationsPath := resolveMigrationsPath()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Error setting dialect: %v", err)
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully")
}

func runRollback(steps int64) {
	db, err := openDb()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	migrationsPath := resolveMigrationsPath()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Error setting dialect: %v", err)
	}

	if err := goose.DownTo(db, migrationsPath, steps); err != nil {
		log.Fatalf("Error performing rollback: %v", err)
	}

	fmt.Printf("Rollback of %d migrations performed successfully\n", steps)
}

func runMigrateToSequential(targetSeq int64) {
	db, err := openDb()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	migrationsPath := resolveMigrationsPath()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Error setting dialect: %v", err)
	}

	if err := goose.UpTo(db, migrationsPath, targetSeq); err != nil {
		log.Fatalf("Error migrating to sequential %05d: %v", targetSeq, err)
	}

	fmt.Printf("Successfully migrated to %05d\n", targetSeq)
}

func runRollbackToSequential(targetSeq int64) {
	db, err := openDb()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	migrationsPath := resolveMigrationsPath()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Error setting dialect: %v", err)
	}

	if err := goose.DownTo(db, migrationsPath, targetSeq); err != nil {
		log.Fatalf("Error rolling back to sequential %05d: %v", targetSeq, err)
	}

	fmt.Printf("Successfully rolled back to %05d\n", targetSeq)
}
