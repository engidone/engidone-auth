package seeders

import (
	"time"

	"engidone-auth/internal/signin/domain"
)

// Seeder represents a database seeder interface
type Seeder interface {
	Name() string
	Seed(db interface{}) error
	Dependencies() []string
}

// BaseSeeder provides common functionality for all seeders
type BaseSeeder struct {
	name         string
	dependencies []string
}

func NewBaseSeeder(name string, dependencies ...string) *BaseSeeder {
	return &BaseSeeder{
		name:         name,
		dependencies: dependencies,
	}
}

func (s *BaseSeeder) Name() string {
	return s.name
}

func (s *BaseSeeder) Dependencies() []string {
	return s.dependencies
}

// UserSeederData represents user data structure for YAML seeders
type UserSeederData struct {
	Users []UserSeed `yaml:"users"`
}

// UserSeed represents a single user seed
type UserSeed struct {
	Username string    `yaml:"username"`
	Email    string    `yaml:"email"`
	Password string    `yaml:"password"`
	Roles    []string  `yaml:"roles,omitempty"`
	Active   *bool     `yaml:"active,omitempty"`
	Metadata *Metadata `yaml:"metadata,omitempty"`
}

// Metadata represents additional user metadata
type Metadata struct {
	FirstName string `yaml:"first_name,omitempty"`
	LastName  string `yaml:"last_name,omitempty"`
	Phone     string `yaml:"phone,omitempty"`
	CreatedAt string `yaml:"created_at,omitempty"`
}

// ToDomainUser converts UserSeed to domain.User
func (u *UserSeed) ToDomainUser() *domain.User {
	user := &domain.User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}

	// Set active flag if provided
	if u.Active != nil {
		// You could add an Active field to domain.User if needed
	}

	// Parse created_at if provided
	if u.Metadata != nil && u.Metadata.CreatedAt != "" {
		if createdAt, err := time.Parse(time.RFC3339, u.Metadata.CreatedAt); err == nil {
			user.CreatedAt = createdAt
			user.UpdatedAt = createdAt
		}
	} else {
		now := time.Now()
		user.CreatedAt = now
		user.UpdatedAt = now
	}

	return user
}

// SeederConfig represents configuration for seeders
type SeederConfig struct {
	Seeders []SeederConfigItem `yaml:"seeders"`
	Enabled bool               `yaml:"enabled"`
}

// SeederConfigItem represents individual seeder configuration
type SeederConfigItem struct {
	Name    string   `yaml:"name"`
	Enabled bool     `yaml:"enabled"`
	Order   int      `yaml:"order"`
	Depends []string `yaml:"depends,omitempty"`
}

// SeederResult represents the result of running a seeder
type SeederResult struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	Duration  time.Duration `json:"duration"`
	Records   int           `json:"records"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}