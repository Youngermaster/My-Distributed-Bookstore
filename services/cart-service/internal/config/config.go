package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	HTTPPort string
	Env      string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Cart Settings
	CartTTLHours         int
	MaxItemsPerCart      int
	MaxQuantityPerItem   int
}

func Load() (*Config, error) {
	// Load .env file if exists (ignore error in production)
	_ = godotenv.Load()

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	cartTTL, _ := strconv.Atoi(getEnv("CART_TTL_HOURS", "168")) // 7 days default
	maxItems, _ := strconv.Atoi(getEnv("MAX_ITEMS_PER_CART", "50"))
	maxQty, _ := strconv.Atoi(getEnv("MAX_QUANTITY_PER_ITEM", "99"))

	return &Config{
		HTTPPort: getEnv("HTTP_PORT", "8083"),
		Env:      getEnv("ENV", "development"),
		
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
		
		CartTTLHours:       cartTTL,
		MaxItemsPerCart:    maxItems,
		MaxQuantityPerItem: maxQty,
	}, nil
}

func (c *Config) Print() {
	log.Println("=== Cart Service Configuration ===")
	log.Printf("HTTP Port: %s", c.HTTPPort)
	log.Printf("Environment: %s", c.Env)
	log.Printf("Redis: %s:%s (DB: %d)", c.RedisHost, c.RedisPort, c.RedisDB)
	log.Printf("Cart TTL: %d hours", c.CartTTLHours)
	log.Printf("Max Items: %d, Max Qty: %d", c.MaxItemsPerCart, c.MaxQuantityPerItem)
	log.Println("=================================")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
