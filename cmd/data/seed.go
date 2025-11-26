package main

import (
	"database/sql"
	"engidoneauth/log"
	"flag"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// runSeedApplication orchestrates the entire seed application process
func runSeedApplication() {
	// Parse command line flags
	truncateFlag := flag.Bool("truncate", false, "Truncate tables before seeding data")
	tableFlag := flag.String("table", "", "Specify which table(s) to seed (comma-separated). If not specified, all tables will be processed")
	flag.Parse()

	// Initialize paths and database
	seedsDirPath := getSeedsDirPath()
	db := initializeDatabase()
	defer db.Close()

	// Load available seed files
	seedFiles, err := LoadAllSeedFiles(seedsDirPath)
	if err != nil {
		log.Fatalf("Error loading seed files: %v", err)
	}

	// Extract table names from seed files
	availableTables := extractTableNames(seedFiles)

	// Set up table filtering
	tableFilter := NewTableFilter(*tableFlag, availableTables)

	// Validate specified tables
	validateTableSelection(tableFilter)

	// Handle table truncation if requested
	if *truncateFlag {
		handleTableTruncation(db, tableFilter)
	}

	// Process seed files
	summary := processSeedFiles(db, seedFiles, tableFilter)

	// Print execution summary
	printExecutionSummary(summary, tableFilter)
}

// getSeedsDirPath returns the path to the seeds directory
func getSeedsDirPath() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("Could not get file path")
	}
	return path.Join(filepath.Dir(filename), "seeds")
}

// initializeDatabase creates and returns a database connection
func initializeDatabase() *sql.DB {
	db, err := openDb()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	return db
}

// extractTableNames extracts table names from seed files
func extractTableNames(seedFiles []*SeedFile) []string {
	tableNames := make([]string, 0, len(seedFiles))
	for _, seedFile := range seedFiles {
		tableNames = append(tableNames, seedFile.Data.Name)
	}
	return tableNames
}

// validateTableSelection validates and reports on table selection
func validateTableSelection(tableFilter *TableFilter) {
	if len(tableFilter.SpecifiedTables) > 0 {
		fmt.Printf("🎯 Specified tables: %v\n", tableFilter.SpecifiedTables)

		invalidTables := tableFilter.ValidateSpecifiedTables()
		if len(invalidTables) > 0 {
			fmt.Printf("⚠️  Tables not found: %v\n", invalidTables)
			fmt.Printf("📋 Available tables: %v\n\n", tableFilter.AvailableTables)
			log.Fatalf("Please specify only valid tables")
		}
	} else {
		fmt.Println("📋 No tables specified. All available tables will be processed.")
	}
}

// handleTableTruncation truncates tables if the truncate flag is set
func handleTableTruncation(db *sql.DB, tableFilter *TableFilter) {
	targetTables := tableFilter.GetTargetTables()

	if len(tableFilter.SpecifiedTables) > 0 {
		fmt.Printf("🚀 --truncate flag detected. Only specified tables will be truncated.\n")
	} else {
		fmt.Printf("🚀 --truncate flag detected. All tables will be truncated.\n")
	}

	truncateTables(db, targetTables)
	fmt.Println("")
}

// processSeedFiles processes all seed files and inserts data
func processSeedFiles(db *sql.DB, seedFiles []*SeedFile, tableFilter *TableFilter) *ExecutionSummary {
	processor := NewSeedProcessor()
	summary := &ExecutionSummary{}

	for _, seedFile := range seedFiles {
		// Check if this table should be processed
		if !tableFilter.ShouldProcessTable(seedFile.Data.Name) {
			summary.SkippedCount++
			fmt.Printf("⏭️  Skipping table: %s (not specified)\n", seedFile.Data.Name)
			continue
		}

		// Process the seed file
		result, err := processor.ProcessSeedFile(seedFile)
		if err != nil {
			log.Fatalf("Error processing seed file %s: %v", seedFile.FileName, err)
		}

		// Execute database insertion
		totalInserted := executeSeedInsertion(db, result)

		// Update summary
		summary.ProcessedCount++
		summary.TotalInserted += totalInserted

		// Print progress
		printProgress(result, totalInserted)
	}

	return summary
}

// executeSeedInsertion executes the SQL insert statements for a seed file
func executeSeedInsertion(db *sql.DB, result *ProcessingResult) int {
	fmt.Printf("🔁 Processing file: %s (%d records)\n", result.FileName, result.ItemCount)
	fmt.Printf("📥 Table: %s\n", result.TableName)
	fmt.Printf("🎛️  Records to insert: %d\n", result.ItemCount)

	totalInserted := 0
	for _, values := range result.Values {
		dbResult, err := db.Exec(result.SQL, values...)
		if err != nil {
			log.Fatalf("Error executing seed %v", err)
		}

		if rowsAffected, err := dbResult.RowsAffected(); err == nil {
			totalInserted += int(rowsAffected)
		}
	}

	fmt.Printf("✅ Inserted %d records into table %s from file %s\n",
		totalInserted, result.TableName, result.FileName)
	fmt.Printf("✅ Total inserted: %d records\n\n", totalInserted)

	return totalInserted
}

// printProgress prints progress information for a processed seed file
func printProgress(result *ProcessingResult, totalInserted int) {
	// Progress information is already printed in executeSeedInsertion
}

// ExecutionSummary contains the summary of seed execution
type ExecutionSummary struct {
	ProcessedCount int
	SkippedCount   int
	TotalInserted  int
}

// printExecutionSummary prints a comprehensive summary of the seed execution
func printExecutionSummary(summary *ExecutionSummary, tableFilter *TableFilter) {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("📊 EXECUTION SUMMARY:\n")
	fmt.Printf("📁 Files processed: %d\n", summary.ProcessedCount)
	if summary.SkippedCount > 0 {
		fmt.Printf("⏭️  Files skipped: %d\n", summary.SkippedCount)
	}
	fmt.Printf("🎯 Target tables: %v\n", tableFilter.GetTargetTables())
	fmt.Printf("💾 Total records inserted: %d\n", summary.TotalInserted)
	fmt.Println(strings.Repeat("=", 50))
}

func main() {
	runSeedApplication()
}
