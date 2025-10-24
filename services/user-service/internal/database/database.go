package database

import (
	"fmt"
	"log"

	"github.com/youngermaster/distributed-bookstore/user-service/internal/config"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Set log level
	logLevel := logger.Silent
	if cfg.Env == "development" {
		logLevel = logger.Info
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)

	log.Println("Database connection established")

	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.Address{},
		&domain.Session{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")

	return nil
}

// SeedDefaultRoles creates default roles if they don't exist
func SeedDefaultRoles(db *gorm.DB) error {
	log.Println("Seeding default roles...")

	roles := []domain.Role{
		{
			Name:        "customer",
			Description: "Default customer role",
			Permissions: `["read:own_profile", "write:own_profile", "read:catalog", "create:order"]`,
		},
		{
			Name:        "admin",
			Description: "Administrator role with full access",
			Permissions: `["*"]`,
		},
	}

	for _, role := range roles {
		var existingRole domain.Role
		result := db.Where("name = ?", role.Name).First(&existingRole)
		
		if result.Error != nil {
			// Role doesn't exist, create it
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to create role %s: %v", role.Name, err)
			} else {
				log.Printf("Created role: %s", role.Name)
			}
		} else {
			log.Printf("Role %s already exists", role.Name)
		}
	}

	log.Println("Default roles seeded successfully")

	return nil
}
