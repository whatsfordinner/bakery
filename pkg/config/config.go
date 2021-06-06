package config

import (
	"context"
	"flag"
)

// Config represents all the configuration information for the baker and reception
type Config struct {
	DBHost         *string
	RabbitHost     *string
	RabbitUsername *string
	RabbitPassword *string
}

// GetConfig reads in config values from their appropriate command line flags
func GetConfig(ctx context.Context) *Config {
	c := new(Config)
	c.DBHost = flag.String("dbhost", "127.0.0.1:6379", "connection string for Redis DB")
	c.RabbitHost = flag.String("rabbithost", "amqp://127.0.0.1.5672", "host for RabbitMQ")
	c.RabbitUsername = flag.String("rabbituser", "guest", "username for connecting to RabbitMQ")
	c.RabbitPassword = flag.String("rabbitpass", "guest", "password for connecting to RabbitMQ")
	flag.Parse()

	return c
}
