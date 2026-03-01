package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration values.
type Config struct {
	// REST
	RESTPort  int    // Port for the REST API
	GinMode   string // Mode for the Gin framework (e.g., release, debug, test)
	JWTSecret string // Secret key for JWT signing
	JWTIssuer string // Issuer claim for JWTs

	// MongoDB
	DBHost     string // Hostname or IP address for the database
	DBPort     int    // Port number for the database
	DBUser     string // Username for the database
	DBPassword string // Password for the database
	DBName     string // Name of the database

	// Redis
	RedisHost string // Hostname or IP address for the Redis server
	RedisPort int    // Port number for the Redis server

	// Matchmaking
	MaxPlayer        int32 // Maximum number of players allowed in a game
	RankTolerance    int32 // Tolerance for player rank difference during matchmaking
	LatencyTolerance int32 // Tolerance for latency (in milliseconds) during matchmaking

	// UDP game server
	UdpPort                int // Port for the UDP socket
	UDPBufferSize          int // Size of the buffer for incoming UDP packets (in bytes)
	UDPHeartbeatExpiration int // Expiration time for UDP heartbeat (in milliseconds)
}

// Envs holds the application's configuration loaded from environment variables.
var Envs = initConfig()

// initConfig initializes and returns the application configuration.
// It loads environment variables from a .env file.
func initConfig() Config {
	// Load .env file if available
	if err := godotenv.Load(); err != nil {
		log.Printf("[APP] [INFO] .env file not found or could not be loaded: %v", err)
	}

	// Populate the Config struct with required environment variables
	return Config{
		// REST
		RESTPort:  mustGetEnvAsInt("REST_PORT"),
		GinMode:   getEnvWithDefault("GIN_MODE", "release"),
		JWTSecret: mustGetEnv("JWT_SECRET"),
		JWTIssuer: mustGetEnv("JWT_ISSUER"),

		// MongoDB
		DBHost:     mustGetEnv("DB_HOST"),
		DBPort:     mustGetEnvAsInt("DB_PORT"),
		DBUser:     mustGetEnv("DB_USER"),
		DBPassword: mustGetEnv("DB_PASS"),
		DBName:     mustGetEnv("DB_NAME"),

		// Redis
		RedisHost: mustGetEnv("REDIS_HOST"),
		RedisPort: mustGetEnvAsInt("REDIS_PORT"),

		// Matchmaking
		MaxPlayer:        int32(mustGetEnvAsInt("MAX_PLAYER")),
		RankTolerance:    int32(mustGetEnvAsInt("RANK_TOLERANCE")),
		LatencyTolerance: int32(mustGetEnvAsInt("LATENCY_TOLERANCE")),

		// UDP game server
		UdpPort:                mustGetEnvAsInt("UDP_PORT"),
		UDPBufferSize:          mustGetEnvAsInt("UDP_BUFFER_SIZE"),
		UDPHeartbeatExpiration: mustGetEnvAsInt("UDP_HEARTBEAT_EXPIRATION"),
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

// mustGetEnvAsInt retrieves the value of an environment variable as an integer or logs a fatal error if not set or cannot be parsed.
func mustGetEnvAsInt(key string) int {
	valueStr := mustGetEnv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("[APP] [FATAL] Environment variable %s must be an integer: %v", key, err)
	}
	return value
}

// getEnvWithDefault retrieves the value of an environment variable or returns a default value if not set.
func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
