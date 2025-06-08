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
	cfg := Config{
		ConnectionString: "localhost:6379",
		ServerPort:       3000,
	}

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get working dir %w", err)
	}

	dbPath := filepath.Join(basePath, "db", "app.db")

	cfg.ConnectionString = dbPath

	if serverPort, exists := os.LookupEnv(("SERVER_PORT")); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint64(port)
		}
	}

	return cfg
}
