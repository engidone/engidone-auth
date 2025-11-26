package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

// SeedItems represents a collection of seed data items
type SeedItems []map[string]any

// Func represents function mappings for seed processing
type Func map[string]string

// Seed represents the complete seed configuration from YAML
type Seed struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Items       SeedItems `yaml:"items"`
}

// SeedFile represents a seed file with its path and parsed data
type SeedFile struct {
	Path     string
	FileName string
	Data     Seed
}

// LoadSeedFile loads and parses a single seed file
func LoadSeedFile(filePath string) (*SeedFile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	var seedsData Seed
	err = yaml.Unmarshal(data, &seedsData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML from %s: %w", filePath, err)
	}

	return &SeedFile{
		Path:     filePath,
		FileName: path.Base(filePath),
		Data:     seedsData,
	}, nil
}

// LoadAllSeedFiles loads all seed files from a directory
func LoadAllSeedFiles(seedsDirPath string) ([]*SeedFile, error) {
	entries, err := os.ReadDir(seedsDirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading seeds directory %s: %w", seedsDirPath, err)
	}

	var seedFiles []*SeedFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := path.Join(seedsDirPath, entry.Name())
		seedFile, err := LoadSeedFile(filePath)
		if err != nil {
			// Skip invalid files but continue processing others
			continue
		}

		seedFiles = append(seedFiles, seedFile)
	}

	return seedFiles, nil
}

// GetAvailableTableNames extracts table names from seed files
func GetAvailableTableNames(seedsDirPath string, entries []os.DirEntry) ([]string, error) {
	var availableTables []string

	for _, file := range entries {
		if file.IsDir() {
			continue
		}

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

// GetFieldsAndSequences determines the field order and parameter placeholders for SQL
func GetFieldsAndSequences(data SeedItems) (fields string, sequences string) {
	if len(data) == 0 {
		return "", ""
	}

	var keys []string
	var seqSlice []string

	// Define field order for specific table types to ensure consistency
	if fieldOrder := getTableFieldOrder(data[0]); len(fieldOrder) > 0 {
		keys = fieldOrder
	} else {
		// Use the order from the first item if table type not recognized
		for key := range data[0] {
			keys = append(keys, key)
		}
	}

	// Create parameter placeholders
	for i := 1; i <= len(keys); i++ {
		seqSlice = append(seqSlice, fmt.Sprintf("$%d", i))
	}

	return strings.Join(keys, ", "), strings.Join(seqSlice, ", ")
}

// getTableFieldOrder returns the expected field order for known table types
func getTableFieldOrder(firstItem map[string]any) []string {
	// Check for users table
	if _, hasUserId := firstItem["user_id"]; hasUserId {
		if _, hasPassword := firstItem["password"]; hasPassword {
			return []string{"user_id", "password"}
		}
	}

	// Check for recovery_codes table
	if _, hasUserId := firstItem["user_id"]; hasUserId {
		if _, hasCode := firstItem["code"]; hasCode {
			if _, hasIsValid := firstItem["is_valid"]; hasIsValid {
				if _, hasExpiresAt := firstItem["expires_at"]; hasExpiresAt {
					return []string{"user_id", "code", "is_valid", "expires_at"}
				}
			}
		}
	}

	// Check for refresh_tokens table
	if _, hasUserId := firstItem["user_id"]; hasUserId {
		if _, hasTokenId := firstItem["token_id"]; hasTokenId {
			if _, hasExpiresAt := firstItem["expires_at"]; hasExpiresAt {
				return []string{"user_id", "token_id", "expires_at"}
			}
		}
	}

	return nil // No specific order found
}