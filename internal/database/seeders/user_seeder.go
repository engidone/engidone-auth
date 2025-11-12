package seeders

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"engidone-auth/internal/signin/infrastructure"
)

// UserSeeder implements the Seeder interface for users
type UserSeeder struct {
	*BaseSeeder
	dataPath string
}

// NewUserSeeder creates a new user seeder with default data
func NewUserSeeder() *UserSeeder {
	return &UserSeeder{
		BaseSeeder: NewBaseSeeder("users"),
		dataPath:   "",
	}
}

// NewUserSeederFromConfig creates a user seeder from YAML configuration
func NewUserSeederFromConfig(dataPath string, dependencies []string) *UserSeeder {
	return &UserSeeder{
		BaseSeeder: NewBaseSeeder("users", dependencies...),
		dataPath:   dataPath,
	}
}

// Seed seeds users into the database
func (s *UserSeeder) Seed(db interface{}) error {
	sqlDB, ok := db.(*sql.DB)
	if !ok {
		return fmt.Errorf("expected *sql.DB, got %T", db)
	}

	// Load user data
	userData, err := s.loadUserData()
	if err != nil {
		return fmt.Errorf("failed to load user data: %w", err)
	}

	// Create repository
	repo := infrastructure.NewSQLUserRepository(sqlDB)

	// Seed users
	for _, userSeed := range userData.Users {
		domainUser := userSeed.ToDomainUser()

		// Check if user already exists
		existingUser, err := repo.FindByUsername(domainUser.Username)
		if err == nil && existingUser != nil {
			// User exists, skip or update based on your preference
			continue
		}

		// Create new user
		if err := repo.Create(domainUser); err != nil {
			return fmt.Errorf("failed to create user %s: %w", domainUser.Username, err)
		}
	}

	return nil
}

// loadUserData loads user data from YAML file
func (s *UserSeeder) loadUserData() (*UserSeederData, error) {
	var dataPath string

	if s.dataPath != "" {
		// Use custom data path
		dataPath = filepath.Join(s.dataPath, "data.yml")
	} else {
		// Use default data path
		dataPath = "internal/database/seeders/users/data.yml"
	}

	// If custom data file doesn't exist, use default data
	if _, err := ioutil.ReadFile(dataPath); err != nil {
		return s.getDefaultUserData(), nil
	}

	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var userData UserSeederData
	if err := yaml.Unmarshal(data, &userData); err != nil {
		return nil, err
	}

	return &userData, nil
}

// getDefaultUserData returns default user data
func (s *UserSeeder) getDefaultUserData() *UserSeederData {
	return &UserSeederData{
		Users: []UserSeed{
			{
				Username: "admin",
				Email:    "admin@example.com",
				Password: "admin123",
				Roles:    []string{"admin", "superuser"},
				Active:   boolPtr(true),
				Metadata: &Metadata{
					FirstName: "Administrator",
					LastName:  "User",
					CreatedAt: "2024-01-01T00:00:00Z",
				},
			},
			{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "test123",
				Roles:    []string{"user"},
				Active:   boolPtr(true),
				Metadata: &Metadata{
					FirstName: "Test",
					LastName:  "User",
					CreatedAt: "2024-01-01T00:00:00Z",
				},
			},
			{
				Username: "john",
				Email:    "john@example.com",
				Password: "john123",
				Roles:    []string{"user"},
				Active:   boolPtr(true),
				Metadata: &Metadata{
					FirstName: "John",
					LastName:  "Doe",
					Phone:     "+1234567890",
					CreatedAt: "2024-01-01T00:00:00Z",
				},
			},
			{
				Username: "jane",
				Email:    "jane@example.com",
				Password: "jane123",
				Roles:    []string{"user", "moderator"},
				Active:   boolPtr(true),
				Metadata: &Metadata{
					FirstName: "Jane",
					LastName:  "Smith",
					Phone:     "+0987654321",
					CreatedAt: "2024-01-01T00:00:00Z",
				},
			},
			{
				Username: "developer",
				Email:    "dev@example.com",
				Password: "dev123",
				Roles:    []string{"user", "developer"},
				Active:   boolPtr(true),
				Metadata: &Metadata{
					FirstName: "Developer",
					LastName:  "User",
					CreatedAt: "2024-01-01T00:00:00Z",
				},
			},
		},
	}
}

// boolPtr returns a pointer to a bool
func boolPtr(b bool) *bool {
	return &b
}

// hashPassword generates a SHA-256 hash of the password
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}
