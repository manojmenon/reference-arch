package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds runtime settings loaded from the environment.
type Config struct {
	Port            string
	DatabaseURL     string
	LogLevel        string
	ShutdownTimeout time.Duration
	DBMaxRetries    int
	DBRetryBackoff  time.Duration
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getenvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

// Load reads configuration from environment variables.
func Load() Config {
	return Config{
		Port:            getenv("PORT", "8080"),
		DatabaseURL:     getenv("DATABASE_URL", "postgres://appuser:apppass@localhost:5432/appdb?sslmode=disable"),
		LogLevel:        getenv("LOG_LEVEL", "info"),
		ShutdownTimeout: getenvDuration("SHUTDOWN_TIMEOUT", 15*time.Second),
		DBMaxRetries:    getenvInt("DB_MAX_RETRIES", 5),
		DBRetryBackoff:  getenvDuration("DB_RETRY_BACKOFF", 200*time.Millisecond),
	}
}
