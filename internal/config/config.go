package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
}

func Load() (Config, error) {
	databaseURL := os.Getenv("SPORTS_DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("SPORTS_DATABASE_URL is not set")
	}

	httpAddr := os.Getenv("SPORTS_HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	return Config{
		DatabaseURL: databaseURL,
		HTTPAddr:    httpAddr,
	}, nil
}
