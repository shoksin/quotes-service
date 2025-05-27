package configs

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}
type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

var (
	once sync.Once
	cfg  *Config
)

func Load() *Config {
	once.Do(func() {
		cfg = &Config{
			Server: ServerConfig{
				Port: getEnv("SERVER_PORT", "8080"),
			},
			Database: DatabaseConfig{
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnvInt("DB_PORT", 5432),
				Name:     getEnv("DB_NAME", "quotes_db"),
				User:     getEnv("DB_USER", "quotes_user"),
				Password: getEnv("DB_PASSWORD", "quotes_password"),
			},
		}
	})
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	valueString := os.Getenv(key)
	if value, err := strconv.Atoi(valueString); err == nil {
		return value
	}
	return defaultValue
}
