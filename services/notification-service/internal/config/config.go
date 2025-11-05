package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort              string
	Env                   string
	ServiceName           string
	LogLevel              string
	ShutdownTimeout       time.Duration
	RabbitMQURL           string
	RabbitMQExchange      string
	RabbitMQQueue         string
	RabbitRoutingKeys     []string
	RabbitMQPrefetchCount int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPPort:              getEnv("HTTP_PORT", "8089"),
		Env:                   getEnv("ENV", "development"),
		ServiceName:           getEnv("SERVICE_NAME", "notification-service"),
		LogLevel:              getEnv("LOG_LEVEL", "info"),
		ShutdownTimeout:       getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		RabbitMQURL:           getEnv("RABBITMQ_URL", "amqp://bookstore:dev_password@rabbitmq:5672/"),
		RabbitMQExchange:      getEnv("RABBITMQ_EXCHANGE", "bookstore.events"),
		RabbitMQQueue:         getEnv("RABBITMQ_QUEUE", "notification-service"),
		RabbitRoutingKeys:     getSliceEnv("RABBITMQ_ROUTING_KEYS", []string{"user.*", "wishlist.*", "cart.*", "order.*"}),
		RabbitMQPrefetchCount: getIntEnv("RABBITMQ_PREFETCH_COUNT", 10),
	}

	return cfg, nil
}

func (c *Config) Print() {
	log.Println("=== Notification Service Configuration ===")
	log.Printf("Environment: %s", c.Env)
	log.Printf("HTTP Port: %s", c.HTTPPort)
	log.Printf("Service Name: %s", c.ServiceName)
	log.Printf("RabbitMQ URL: %s", c.RabbitMQURL)
	log.Printf("RabbitMQ Exchange: %s", c.RabbitMQExchange)
	log.Printf("RabbitMQ Queue: %s", c.RabbitMQQueue)
	log.Printf("Routing Keys: %s", strings.Join(c.RabbitRoutingKeys, ", "))
	log.Printf("Prefetch Count: %d", c.RabbitMQPrefetchCount)
	log.Println("==========================================")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		items := strings.Split(value, ",")
		result := make([]string, 0, len(items))
		for _, item := range items {
			trimmed := strings.TrimSpace(item)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}
