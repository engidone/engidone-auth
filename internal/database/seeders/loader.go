package seeders

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

	"gopkg.in/yaml.v3"
)

// SeederManager manages all database seeders
type SeederManager struct {
	db           *sql.DB
	seedersPath  string
	seeders      map[string]Seeder
	config       *SeederConfig
	logger       *log.Logger
}

// NewSeederManager creates a new seeder manager
func NewSeederManager(db *sql.DB, seedersPath string, logger *log.Logger) *SeederManager {
	if logger == nil {
		logger = log.New(os.Stdout, "[SEEDERS] ", log.LstdFlags)
	}

	return &SeederManager{
		db:          db,
		seedersPath: seedersPath,
		seeders:     make(map[string]Seeder),
		logger:      logger,
	}
}

// LoadSeeders loads all available seeders from the seeders directory
func (sm *SeederManager) LoadSeeders() error {
	// Load seeder configuration
	if err := sm.loadConfig(); err != nil {
		return fmt.Errorf("failed to load seeder config: %w", err)
	}

	// Register built-in seeders
	sm.registerBuiltInSeeders()

	// Load custom seeders from directory
	if err := sm.loadCustomSeeders(); err != nil {
		return fmt.Errorf("failed to load custom seeders: %w", err)
	}

	sm.logger.Printf("Loaded %d seeders", len(sm.seeders))
	return nil
}

// RunSeeders runs all configured seeders
func (sm *SeederManager) RunSeeders() ([]SeederResult, error) {
	if sm.config == nil || !sm.config.Enabled {
		sm.logger.Println("Seeders are disabled")
		return nil, nil
	}

	var results []SeederResult

	// Sort seeders by order
	sort.Slice(sm.config.Seeders, func(i, j int) bool {
		return sm.config.Seeders[i].Order < sm.config.Seeders[j].Order
	})

	// Run each enabled seeder
	for _, seederConfig := range sm.config.Seeders {
		if !seederConfig.Enabled {
			continue
		}

		result := sm.runSeeder(seederConfig.Name)
		results = append(results, result)
	}

	return results, nil
}

// RunSeeder runs a specific seeder by name
func (sm *SeederManager) RunSeeder(name string) SeederResult {
	return sm.runSeeder(name)
}

// runSeeder runs a single seeder
func (sm *SeederManager) runSeeder(name string) SeederResult {
	result := SeederResult{
		Name:      name,
		Timestamp: time.Now(),
	}

	start := time.Now()

	seeder, exists := sm.seeders[name]
	if !exists {
		result.Status = "failed"
		result.Error = fmt.Sprintf("seeder '%s' not found", name)
		result.Duration = time.Since(start)
		return result
	}

	sm.logger.Printf("Running seeder: %s", name)

	// Check dependencies
	for _, dep := range seeder.Dependencies() {
		if _, exists := sm.seeders[dep]; !exists {
			result.Status = "failed"
			result.Error = fmt.Sprintf("dependency '%s' not found for seeder '%s'", dep, name)
			result.Duration = time.Since(start)
			return result
		}
	}

	// Execute seeder
	if err := seeder.Seed(sm.db); err != nil {
		result.Status = "failed"
		result.Error = err.Error()
		result.Duration = time.Since(start)
		sm.logger.Printf("Seeder %s failed: %v", name, err)
		return result
	}

	result.Status = "success"
	result.Duration = time.Since(start)
	result.Records = 1 // This could be enhanced to return actual record count
	sm.logger.Printf("Seeder %s completed successfully in %v", name, result.Duration)

	return result
}

// loadConfig loads seeder configuration from YAML file
func (sm *SeederManager) loadConfig() error {
	configPath := filepath.Join(sm.seedersPath, "seeders.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config if it doesn't exist
		sm.config = &SeederConfig{
			Enabled: true,
			Seeders: []SeederConfigItem{
				{Name: "users", Enabled: true, Order: 1},
			},
		}
		return sm.saveDefaultConfig(configPath)
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &sm.config)
}

// saveDefaultConfig saves a default seeder configuration
func (sm *SeederManager) saveDefaultConfig(configPath string) error {
	data, err := yaml.Marshal(sm.config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, data, 0644)
}

// registerBuiltInSeeders registers built-in seeders
func (sm *SeederManager) registerBuiltInSeeders() {
	// Register user seeder
	userSeeder := NewUserSeeder()
	sm.seeders[userSeeder.Name()] = userSeeder
}

// loadCustomSeeders loads custom seeders from the seeders directory
func (sm *SeederManager) loadCustomSeeders() error {
	entries, err := ioutil.ReadDir(sm.seedersPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			if err := sm.loadSeederFromDirectory(entry.Name()); err != nil {
				sm.logger.Printf("Failed to load seeder from %s: %v", entry.Name(), err)
			}
		}
	}

	return nil
}

// loadSeederFromDirectory loads a seeder from a directory
func (sm *SeederManager) loadSeederFromDirectory(name string) error {
	seederDir := filepath.Join(sm.seedersPath, name)
	yamlFile := filepath.Join(seederDir, "seeder.yml")

	if _, err := os.Stat(yamlFile); os.IsNotExist(err) {
		return nil // Skip if no seeder.yml file exists
	}

	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return err
	}

	var config struct {
		Type         string `yaml:"type"`
		Dependencies []string `yaml:"dependencies"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	// Create seeder based on type
	switch config.Type {
	case "users":
		seeder := NewUserSeederFromConfig(seederDir, config.Dependencies)
		sm.seeders[seeder.Name()] = seeder
	default:
		return fmt.Errorf("unknown seeder type: %s", config.Type)
	}

	return nil
}