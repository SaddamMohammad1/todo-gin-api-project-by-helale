package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds the application's configuration values.
// These values are usually loaded from a .env file or system environment variables.
type Config struct {
	DatabaseURL string // Example: "postgres://user:pass@localhost:5432/dbname"
	Port        string // Example: "8080" - server port number
	JWTSecret   string
}

// Load function loads environment variables and returns a Config object.
// It first tries to load values from a .env file.
// If .env is not found, it continues using system environment variables.
// Note - Without godotenv.Load() → Go cannot read .env file
/*
	Q. What does godotenv.Load() do?
	Ans: It loads the .env file into system environment variables at runtime
		This function:
			Reads the .env file
			Takes each line (e.g., DATABASE_URL=xxx)
			Stores them into the process environment
			After that, os.Getenv("DATABASE_URL") can read it
*/
func Load() (*Config, error) {
	var err error = godotenv.Load()

	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	var config *Config = &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("Port"),
		JWTSecret:   os.Getenv("JWT_Secret"),
	}

	return config, nil
}
