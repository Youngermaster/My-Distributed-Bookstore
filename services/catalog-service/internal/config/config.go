package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	// Server
	HTTPPort string
	GRPCPort string
	Env      string

	// Database
	DatabaseURL      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  string

	// Redis
	RedisURL      string
	RedisPassword string

	// RabbitMQ
	RabbitMQURL string

	// Observability
	JaegerAgentHost string
	JaegerAgentPort string
	LogLevel        string
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTPPort: getEnv("HTTP_PORT", "8081"),
		GRPCPort: getEnv("GRPC_PORT", "50051"),
		Env:      getEnv("ENV", "development"),

		// Database
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "5432"),
		DatabaseUser:     getEnv("DATABASE_USER", "bookstore"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "password"),
		DatabaseName:     getEnv("DATABASE_NAME", "catalog_db"),
		DatabaseSSLMode:  getEnv("DATABASE_SSLMODE", "disable"),

		// Redis
		RedisURL:      getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// RabbitMQ
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://bookstore:password@localhost:5672/"),

		// Observability
		JaegerAgentHost: getEnv("JAEGER_AGENT_HOST", "localhost"),
		JaegerAgentPort: getEnv("JAEGER_AGENT_PORT", "6831"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

func (c *Config) GetDatabaseDSN() string {
	// If DATABASE_URL is provided, use it directly
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	// Otherwise, construct from individual parts
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseName,
		c.DatabaseSSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func (c *Config) Print() {
	log.Println("=== Catalog Service Configuration ===")
	log.Printf("Environment: %s", c.Env)
	log.Printf("HTTP Port: %s", c.HTTPPort)
	log.Printf("gRPC Port: %s", c.GRPCPort)
	log.Printf("Database: %s@%s:%s/%s", c.DatabaseUser, c.DatabaseHost, c.DatabasePort, c.DatabaseName)
	log.Printf("Redis: %s", c.RedisURL)
	log.Printf("RabbitMQ: Connected")
	log.Printf("Jaeger: %s:%s", c.JaegerAgentHost, c.JaegerAgentPort)
	log.Println("====================================")
}
