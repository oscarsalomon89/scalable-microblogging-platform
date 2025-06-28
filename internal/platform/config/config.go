package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/oscarsalomon89/go-hexagonal/internal/platform/environment"
)

const (
	scopeEnv               = "GO_ENVIRONMENT"
	localScope             = "local"
	partialProductionScope = "prod"

	AppPathEnv = "APP_PATH"
)

type (
	Configuration struct {
		Scope    string
		Database Database
	}

	Database struct {
		Host               string
		Name               string
		Port               int
		User               string
		Password           string
		Sslmode            string
		MaxConnections     int
		MaxIdleConnections int
		MaxIdleTime        int
		MaxLifeTime        int
		ConnectionTimeout  int
		UseReplica         bool
		BatchSize          int
	}
)

func NewConfig() (Configuration, error) {
	env := environment.GetFromString(os.Getenv(scopeEnv))

	if env == environment.Local {
		if err := godotenv.Load(os.Getenv(AppPathEnv) + "/.env"); err != nil {
			return Configuration{}, err
		}
	}

	return Configuration{
		Scope: getEnv(scopeEnv, localScope),
		Database: Database{
			Name:               getEnv("DB_NAME", "testlocal"),
			Host:               getEnv("DB_HOST", "localhost"),
			Port:               getEnvInt("DB_PORT", 5432),
			User:               getEnv("DB_USER", "postgres"),
			Password:           getEnv("DB_PASS", ""),
			Sslmode:            getEnv("SSL_MODE", "disable"),
			MaxConnections:     getEnvInt("DB_MAX_CONNECTIONS", 100),
			MaxIdleConnections: getEnvInt("DB_MAX_IDLE_CONNECTIONS", 100),
			MaxIdleTime:        getEnvInt("DB_MAX_IDLE_TIME", 600),
			MaxLifeTime:        getEnvInt("DB_MAX_LIFE_TIME", 600),
			ConnectionTimeout:  getEnvInt("DB_CONNECTION_TIMEOUT", 200),
			UseReplica:         getEnvBool("DB_USE_REPLICA", true),
			BatchSize:          getEnvInt("DB_BATCH_SIZE", 200),
		},
	}, nil
}

func getEnv(key string, defaultValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, found := os.LookupEnv(key); found {
		i, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return i
	}

	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, found := os.LookupEnv(key); found {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return b
	}

	return defaultValue
}

func IsLocalScope() bool {
	return getEnv(scopeEnv, localScope) == localScope
}

func IsProductionScope() bool {
	return strings.Contains(getEnv(scopeEnv, localScope), partialProductionScope)
}
