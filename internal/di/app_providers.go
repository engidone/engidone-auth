package di

import (
	"os"

	"github.com/go-kit/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// LoggerModule provides application-level logger
var LoggerModule = fx.Options(
	fx.Provide(
		NewZapLogger,
		NewGoKitLogger,
	),
)

// NewZapLogger creates a new zap logger
func NewZapLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	return config.Build()
}

// NewGoKitLogger creates a new go-kit logger
func NewGoKitLogger() log.Logger {
	return log.NewLogfmtLogger(os.Stdout)
}

// AppConfig holds application configuration
type AppConfig struct {
	ServerPort string
	JWTSecret  string
}

// NewAppConfig creates application configuration
func NewAppConfig() *AppConfig {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "9000"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key"
	}

	return &AppConfig{
		ServerPort: port,
		JWTSecret:  secret,
	}
}

// ConfigModule provides application configuration
var ConfigModule = fx.Options(
	fx.Provide(NewAppConfig),
)
