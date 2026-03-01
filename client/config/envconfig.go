package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration values.
type Config struct {
	ServerAddr  string
	LoginUri    string
	RegisterUri string
	MatchUri    string
}

// Envs holds the application's configuration loaded from environment variables.
var Envs = initConfig()

// initConfig initializes and returns the application configuration.
// It loads environment variables from a .env file.
func initConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("[APP] [INFO] .env file not found or could not be loaded: %v", err)
	}

	// Populate the Config struct with required environment variables
	return Config{
		ServerAddr:  mustGetEnv("SERVER_ADDR"),
		LoginUri:    mustGetEnv("LOGIN_URI"),
		RegisterUri: mustGetEnv("REGISTER_URI"),
		MatchUri:    mustGetEnv("MATCH_URI"),
	}
}

// mustGetEnv retrieves the value of an environment variable or logs a fatal error if not set.
func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("[APP] [FATAL] Environment variable %s is not set", key)
	}
	return value
}
