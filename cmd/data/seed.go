package main

import (
	"database/sql"
	"engidoneauth/internal/config"
	"engidoneauth/log"
	"engidoneauth/util/crypto"
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
		log.Fatal("No se pudo obtener la ruta del archivo")
	}

	return filepath.Dir(filename)
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
	dirPath := getDirPath()
	seedsDirPath := path.Join(dirPath, "seeds")
	entries, err := os.ReadDir(seedsDirPath)

	if err != nil {
		log.Fatalf("Error scanning seeds: %v", err)
	}

	db, _ := openDb()

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
}
