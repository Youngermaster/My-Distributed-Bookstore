package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	HTTPPort       string
	GRPCPort       string
	Env            string

	// Database
	DatabaseURL    string
	DBMaxIdleConns int
	DBMaxOpenConns int

	// JWT
	JWTSecret     string
	JWTExpiration time.Duration

	// Service URLs (gRPC)
	UserServiceURL          string
	CatalogServiceURL       string
	OrderServiceURL         string
	CartServiceURL          string
	InventoryServiceURL     string
	RecommendationServiceURL string

	// CORS
	CORSAllowOrigins string

	// Logging
	LogLevel string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{
		HTTPPort:       getEnv("HTTP_PORT", "8090"),
		GRPCPort:       getEnv("GRPC_PORT", "50060"),
		Env:            getEnv("ENV", "development"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		DBMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		DBMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpiration:  getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),

		// Service URLs
		UserServiceURL:          getEnv("USER_SERVICE_URL", "user-service:50052"),
		CatalogServiceURL:       getEnv("CATALOG_SERVICE_URL", "catalog-service:50051"),
		OrderServiceURL:         getEnv("ORDER_SERVICE_URL", "order-service:50054"),
		CartServiceURL:          getEnv("CART_SERVICE_URL", "cart-service:50053"),
		InventoryServiceURL:     getEnv("INVENTORY_SERVICE_URL", "inventory-service:50056"),
		RecommendationServiceURL: getEnv("RECOMMENDATION_SERVICE_URL", "recommendation-service:8089"),

		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if config.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		return intValue
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}
	return defaultValue
}
