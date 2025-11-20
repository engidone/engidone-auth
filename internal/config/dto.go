package config

type Paths struct {
	Config string
	App    string
}

type AppConfig struct {
	Application struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Port    string `yaml:"port"`
	} `yaml:"application"`
	Database struct {
		DSN      string `yaml:"dsn"`
		Engine   string `yaml:"engine"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		SSLMode  string `yaml:"ssl_mode"`
		DBName   string `yaml:"db_name"`
	} `yaml:"database"`
}
