package main

import (
	"database/sql"
	"engidoneauth/internal/config"
	"engidoneauth/log"
	"engidoneauth/util/crypto"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type SeedItems []map[string]any
type Func map[string]string
type Seed struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Items       SeedItems `yaml:"items"`
}

func handleUUIDString() string {
	return uuid.NewString()
}

func handleUUIDParse(id string) string {
	return uuid.MustParse(id).String()
}

func handleHashPassword(password string) string {
	return crypto.HashPassword(password)
}

func handleAddMinutesToNow(minutes int) time.Time {
	return time.Now().Add(time.Duration(minutes) * time.Minute)
}

func toString(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func getConfigPath(args ...string) string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "..", filepath.Join(args...))
}

func openDb() (*sql.DB, error) {
	appConfig, err := config.LoadFile[config.AppConfig](getConfigPath("config", "app.yaml"))

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	database := appConfig.Database
	dsn := database.DSN
	dsn = strings.ReplaceAll(dsn, "{engine}", database.Engine)
	dsn = strings.ReplaceAll(dsn, "{user}", database.Username)
	dsn = strings.ReplaceAll(dsn, "{password}", database.Password)
	dsn = strings.ReplaceAll(dsn, "{host}", database.Host)
	dsn = strings.ReplaceAll(dsn, "{port}", database.Port)
	dsn = strings.ReplaceAll(dsn, "{db_name}", database.DBName)
	dsn = strings.ReplaceAll(dsn, "{ssl_mode}", database.SSLMode)

	db, err := sql.Open(database.Engine, dsn)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	return db, nil
}

func getDirPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get file path")
	}

	return filepath.Dir(filename)
}

func truncateTables(db *sql.DB, tableNames []string) {
	fmt.Println("🧹 Iniciando truncado de tablas...")

	for _, tableName := range tableNames {
		// Desactivar restricciones de clave externa temporalmente
		disableFKQuery := "SET session_replication_role = replica;"
		_, err := db.Exec(disableFKQuery)
		if err != nil {
			log.Fatalf("Error desactivando restricciones de clave externa: %v", err)
		}

		// Contar registros antes de truncar
		var count int
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
		err = db.QueryRow(countQuery).Scan(&count)
		if err != nil {
			fmt.Printf("⚠️  No se pudo contar registros en tabla %s: %v\n", tableName, err)
		} else {
			fmt.Printf("📊 Tabla %s: %d registros encontrados\n", tableName, count)
		}

		// Truncar tabla
		truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName)
		_, err = db.Exec(truncateQuery)
		if err != nil {
			log.Fatalf("Error truncando tabla %s: %v", tableName, err)
		}

		fmt.Printf("✅ Tabla %s truncada exitosamente\n", tableName)
	}

	// Restaurar restricciones de clave externa
	enableFKQuery := "SET session_replication_role = DEFAULT;"
	_, err := db.Exec(enableFKQuery)
	if err != nil {
		log.Fatalf("Error restaurando restricciones de clave externa: %v", err)
	}

	fmt.Println("🎉 Truncado de tablas completado")
}

func getTableNamesFromFlags(tableFlag string) []string {
	if tableFlag == "" {
		return []string{}
	}

	// Split comma-separated table names and trim whitespace
	tables := strings.Split(tableFlag, ",")
	var result []string
	for _, table := range tables {
		trimmed := strings.TrimSpace(table)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func shouldProcessTable(tableName string, specifiedTables []string) bool {
	// If no tables specified, process all
	if len(specifiedTables) == 0 {
		return true
	}

	// Check if table name is in the specified list
	for _, specified := range specifiedTables {
		if strings.EqualFold(tableName, specified) {
			return true
		}
	}
	return false
}

func getAvailableTables(seedsDirPath string, entries []os.DirEntry) ([]string, error) {
	var availableTables []string

	for _, file := range entries {
		data, err := os.ReadFile(path.Join(seedsDirPath, file.Name()))
		if err != nil {
			continue // Skip files that can't be read
		}

		var seedsData Seed
		err = yaml.Unmarshal(data, &seedsData)
		if err != nil {
			continue // Skip invalid YAML files
		}

		availableTables = append(availableTables, seedsData.Name)
	}

	return availableTables, nil
}

func getFieldsAndSecuences(data SeedItems) (string, string) {
	maplen := len(data[0])
	keys := make([]string, 0, maplen)
	secuences := make([]string, 0, maplen)

	// Define field order for each table type to ensure consistency
	if len(data) > 0 {
		// Try to determine table type based on field names
		if _, hasUserId := data[0]["user_id"]; hasUserId {
			if _, hasPassword := data[0]["password"]; hasPassword {
				keys = []string{"user_id", "password"}
			} else if _, hasCode := data[0]["code"]; hasCode {
				keys = []string{"user_id", "code", "is_valid", "expires_at"}
			}
		}
	}

	// If table not recognized, use original order
	if len(keys) == 0 {
		i := 1
		for k := range data[0] {
			keys = append(keys, k)
			secuences = append(secuences, fmt.Sprintf("$%d", i))
			i++
		}
	} else {
		// Use predefined order
		for i := 1; i <= len(keys); i++ {
			secuences = append(secuences, fmt.Sprintf("$%d", i))
		}
	}

	return strings.Join(keys, ", "), strings.Join(secuences, ", ")
}

func main() {
	// Parse command line flags
	truncateFlag := flag.Bool("truncate", false, "Truncate tables before seeding data")
	tableFlag := flag.String("table", "", "Specify which table(s) to seed (comma-separated). If not specified, all tables will be processed")
	flag.Parse()

	dirPath := getDirPath()
	seedsDirPath := path.Join(dirPath, "seeds")
	entries, err := os.ReadDir(seedsDirPath)

	if err != nil {
		log.Fatalf("Error scanning seeds: %v", err)
	}

	db, _ := openDb()

	// Get specified tables from flag
	specifiedTables := getTableNamesFromFlags(*tableFlag)

	// Get available tables from seed files
	availableTables, err := getAvailableTables(seedsDirPath, entries)
	if err != nil {
		log.Fatalf("Error getting available tables: %v", err)
	}

	// Validate specified tables if any were provided
	if len(specifiedTables) > 0 {
		fmt.Printf("🎯 Tablas especificadas: %v\n", specifiedTables)

		// Check if specified tables exist
		var invalidTables []string
		for _, specified := range specifiedTables {
			found := false
			for _, available := range availableTables {
				if strings.EqualFold(specified, available) {
					found = true
					break
				}
			}
			if !found {
				invalidTables = append(invalidTables, specified)
			}
		}

		if len(invalidTables) > 0 {
			fmt.Printf("⚠️  Tablas no encontradas: %v\n", invalidTables)
			fmt.Printf("📋 Tablas disponibles: %v\n\n", availableTables)
			log.Fatalf("Por favor, especifica solo tablas válidas")
		}
	} else {
		fmt.Println("📋 No se especificaron tablas. Se procesarán todas las tablas disponibles.")
	}

	// If truncate flag is set, collect table names and truncate them
	if *truncateFlag {
		if len(specifiedTables) > 0 {
			fmt.Printf("🚀 Flag --truncate detectado. Se truncarán solo las tablas especificadas.\n")
			truncateTables(db, specifiedTables)
		} else {
			fmt.Printf("🚀 Flag --truncate detectado. Se truncarán todas las tablas.\n")
			truncateTables(db, availableTables)
		}
		fmt.Println("")
	}

	processedCount := 0
	skippedCount := 0

	for _, file := range entries {
		data, err := os.ReadFile(path.Join(seedsDirPath, file.Name()))
		if err != nil {
			log.Fatalf("Error getting data for %s: %s", file.Name(), err.Error())
		}

		var seedsData Seed
		err = yaml.Unmarshal(data, &seedsData)

		if err != nil {
			log.Fatalf("Error unmarshal %s: %s", file.Name(), err.Error())
		}

		// Check if this table should be processed
		if !shouldProcessTable(seedsData.Name, specifiedTables) {
			skippedCount++
			fmt.Printf("⏭️  Omitiendo tabla: %s (no especificada)\n", seedsData.Name)
			continue
		}

		processedCount++
		fields, secuences := getFieldsAndSecuences(seedsData.Items)
		fieldOrder := strings.Split(fields, ", ")
		sentence := "INSERT INTO {tablename} ({fields}) VALUES({secuences})"
		sentence = strings.ReplaceAll(sentence, "{tablename}", seedsData.Name)
		sentence = strings.ReplaceAll(sentence, "{fields}", fields)
		sentence = strings.ReplaceAll(sentence, "{secuences}", secuences)
		// Collect all values for all items in this file
		allValues := make([][]any, 0, len(seedsData.Items))
		for _, item := range seedsData.Items {
			values := make([]any, 0, len(item))
			// Collect values in the same order as fields
			for _, field := range fieldOrder {
				value := item[field]
				splitValues := strings.Split(toString(value), "|")
				if len(splitValues) == 1 {
					switch splitValues[0] {
					case "uuid_string":
						values = append(values, handleUUIDString())
						continue
					case "uuid_parse":
						if len(splitValues) > 1 {
							values = append(values, handleUUIDParse(splitValues[1]))
						} else {
							values = append(values, handleUUIDParse(""))
						}
						continue
					case "hash_password":
						if len(splitValues) > 1 {
							values = append(values, handleHashPassword(splitValues[1]))
						} else {
							values = append(values, handleHashPassword(""))
						}
						continue
					case "add_minutes_to_now":
						if len(splitValues) > 1 {
							if minutes, err := strconv.Atoi(splitValues[1]); err == nil {
								values = append(values, handleAddMinutesToNow(minutes))
							} else {
								values = append(values, handleAddMinutesToNow(0))
							}
						} else {
							values = append(values, handleAddMinutesToNow(0))
						}
						continue
					}
					values = append(values, splitValues[0])
				} else if len(splitValues) == 2 {
					switch splitValues[0] {
					case "uuid_parse":
						values = append(values, handleUUIDParse(splitValues[1]))
						continue
					case "hash_password":
						values = append(values, handleHashPassword(splitValues[1]))
						continue
					case "add_minutes_to_now":
						if minutes, err := strconv.Atoi(splitValues[1]); err == nil {
							values = append(values, handleAddMinutesToNow(minutes))
						} else {
							values = append(values, handleAddMinutesToNow(0))
						}
						continue
					}
					// If not a function, just add the value as-is
					values = append(values, splitValues[1])
				}
			}
			allValues = append(allValues, values)
		}

		// Execute batch insert for all items in this file
		fmt.Printf("🔁 Procesando archivo: %s (%d registros)\n", file.Name(), len(allValues))
		//fmt.Printf("=== Archivo: %s ===\n", file.Name())
		fmt.Printf("📥 Tabla: %s\n", seedsData.Name)
		fmt.Printf("🎛️  Registros a insertar: %d\n", len(allValues))
		//fmt.Printf("SQL: %s\n", sentence)

		// Execute each item in the file (individual INSERTs but grouped by file)
		totalInserted := 0
		for _, values := range allValues {
			//fmt.Printf("Registro %d: %v\n", i+1, values)
			result, err := db.Exec(sentence, values...)
			if err != nil {
				log.Fatalf("Error ejecutando seed %v", err)
			}
			if rowsAffected, err := result.RowsAffected(); err == nil {
				totalInserted += int(rowsAffected)
			}
		}

		fmt.Printf("✅ Se insertaron %d registros en la tabla %s desde el archivo %s\n", totalInserted, seedsData.Name, file.Name())
		fmt.Printf("✅ Total insertado: %d registros\n\n", totalInserted)

	}

	// Print summary
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("📊 RESUMEN DE EJECUCIÓN:\n")
	fmt.Printf("📁 Archivos procesados: %d\n", processedCount)
	if skippedCount > 0 {
		fmt.Printf("⏭️  Archivos omitidos: %d\n", skippedCount)
	}
	fmt.Printf("🎯 Tablas objetivo: %v\n", func() []string {
		if len(specifiedTables) > 0 {
			return specifiedTables
		}
		return availableTables
	}())
	fmt.Println(strings.Repeat("=", 50))
}
