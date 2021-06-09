package config

import (
	"context"
	"os"
)

// Config represents all the configuration information for the baker and reception
type Config struct {
	DBHost         string
	RabbitHost     string
	RabbitUsername string
	RabbitPassword string
}

// GetConfig reads in config values from their appropriate command line flags
func GetConfig(ctx context.Context) *Config {
	c := new(Config)

	// Set sensible defaults (I.e. testing or a dev environment)
	c.DBHost = "127.0.0.1:6379"
	c.RabbitHost = "127.0.0.1:5672"
	c.RabbitUsername = "guest"
	c.RabbitPassword = "guest"

	// Source values from the environment
	val, exists := os.LookupEnv("DB_HOST")
	if exists {
		c.DBHost = val
	}

	val, exists = os.LookupEnv("RABBIT_HOST")
	if exists {
		c.RabbitHost = val
	}

	val, exists = os.LookupEnv("RABBIT_USERNAME")
	if exists {
		c.RabbitUsername = val
	}

	val, exists = os.LookupEnv("RABBIT_PASSWORD")
	if exists {
		c.RabbitPassword = val
	}

	return c
}
