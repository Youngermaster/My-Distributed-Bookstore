package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server
	Port string
	Env  string

	// Service URLs
	CatalogServiceURL        string
	UserServiceURL           string
	CartServiceURL           string
	OrderServiceURL          string
	RecommendationServiceURL string

	// JWT
	JWTSecret     string
	JWTExpiration time.Duration

	// Rate Limiting
	RateLimitEnabled bool
	RateLimitMax     int
	RateLimitWindow  time.Duration

	// Timeouts
	RequestTimeout  time.Duration
	ResponseTimeout time.Duration

	// CORS
	AllowedOrigins string

	// Observability
	JaegerAgentHost string
	JaegerAgentPort string
	LogLevel        string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		// Service URLs
		CatalogServiceURL:        getEnv("CATALOG_SERVICE_URL", "http://localhost:8081"),
		UserServiceURL:           getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		CartServiceURL:           getEnv("CART_SERVICE_URL", "http://localhost:8083"),
		OrderServiceURL:          getEnv("ORDER_SERVICE_URL", "http://localhost:8084"),
		RecommendationServiceURL: getEnv("RECOMMENDATION_SERVICE_URL", "http://localhost:8089"),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),

		// Rate Limiting
		RateLimitEnabled: getBoolEnv("RATE_LIMIT_ENABLED", true),
		RateLimitMax:     getIntEnv("RATE_LIMIT_MAX", 100),
		RateLimitWindow:  getDurationEnv("RATE_LIMIT_WINDOW", 1*time.Minute),

		// Timeouts
		RequestTimeout:  getDurationEnv("REQUEST_TIMEOUT", 30*time.Second),
		ResponseTimeout: getDurationEnv("RESPONSE_TIMEOUT", 30*time.Second),

		// CORS
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),

		// Observability
		JaegerAgentHost: getEnv("JAEGER_AGENT_HOST", "localhost"),
		JaegerAgentPort: getEnv("JAEGER_AGENT_PORT", "6831"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

func (c *Config) Print() {
	fmt.Println("=== API Gateway Configuration ===")
	fmt.Printf("Environment: %s\n", c.Env)
	fmt.Printf("Port: %s\n", c.Port)
	fmt.Println("\n📡 Backend Services:")
	fmt.Printf("  Catalog: %s\n", c.CatalogServiceURL)
	fmt.Printf("  User: %s\n", c.UserServiceURL)
	fmt.Printf("  Cart: %s\n", c.CartServiceURL)
	fmt.Printf("  Order: %s\n", c.OrderServiceURL)
	fmt.Printf("  Recommendation: %s\n", c.RecommendationServiceURL)
	fmt.Printf("\nRate Limiting: Enabled=%v, Max=%d, Window=%v\n", c.RateLimitEnabled, c.RateLimitMax, c.RateLimitWindow)
	fmt.Printf("Request Timeout: %v\n", c.RequestTimeout)
	fmt.Printf("CORS Origins: %s\n", c.AllowedOrigins)
	fmt.Println("=================================")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
