package env

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func InitEnv(filenames ...string) error {
	var err error

	if len(filenames) > 0 {
		err = godotenv.Load(filenames...)
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		return fmt.Errorf("failed to load env file: %w", err)
	}
	return nil
}

func Env(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}

func EnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return i
}

func EnvDuration(key string, defaultValue time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultValue
	}
	return d
}
