package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GanacheURL         string
	DeployerPrivateKey string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		GanacheURL:         getEnv("GANACHE_URL", "http://localhost:8545"),
		DeployerPrivateKey: getEnv("DEPLOYER_PRIVATE_KEY", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
