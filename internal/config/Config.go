package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SQLConnectionString string
}

func New() (*Config, error) {
	err := godotenv.Load()
	
	if err != nil {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	dsn := os.Getenv("SQL_CONNECTION_STRING")
	if dsn == "" {
		return nil, fmt.Errorf("SQL_CONNECTION_STRING is empty")
	}

	return &Config{
		SQLConnectionString: dsn,
	}, nil
}
