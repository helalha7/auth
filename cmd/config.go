package main

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	db struct {
		dsn string
	}
}

func loadConfig() (config, error) {
	if err := godotenv.Load("./.env"); err != nil {
		return config{}, err
	}

	dbDSN, err := getEnv("DB_DSN")
	if err != nil {
		return config{}, err
	}

	cfg := config{
		db: struct{ dsn string }{dsn: dbDSN},
	}

	return cfg, nil
}

func getEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return val, errors.New(key + "not loaded to environment variables")
	}
	return val, nil
}

func getEnvWithDefault(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
