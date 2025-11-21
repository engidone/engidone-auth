package config

type Paths struct {
	Root   string
	Config string
}

type AppConfig struct {
	Application struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Server  struct {
			Port string `yaml:"port"`
		} `yaml:"server"`
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
	Certs struct {
		Public  string `yaml:"public"`
		Private string `yaml:"private"`
	} `yaml:"certs"`
}
