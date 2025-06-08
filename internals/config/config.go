package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	ConnectionString string
	ServerPort       uint64
}

func LoadConfig() Config {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get working dir %w", err)
	}

	dbPath := filepath.Join(basePath, "db", "app.db")
	cfg := Config{
		ConnectionString: dbPath,
		ServerPort:       3000,
	}

	if serverPort, exists := os.LookupEnv(("SERVER_PORT")); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint64(port)
		}
	}

	return cfg
}
