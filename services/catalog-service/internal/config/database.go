package config

import (
	"fmt"
	"log"
	"time"

	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase initializes the database connection
func InitDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseDSN()

	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.Env == "production" {
		gormLogger = logger.Default.LogMode(logger.Error)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")

	return db, nil
}

// AutoMigrate runs automatic migrations for all models
func AutoMigrate(db *gorm.DB) error {
	log.Println("🔄 Running database auto-migrations...")

	err := db.AutoMigrate(
		&domain.Book{},
		&domain.Author{},
		&domain.Category{},
		&domain.Publisher{},
		&domain.BookAuthor{},
		&domain.BookCategory{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// SeedDatabase seeds the database with initial data
func SeedDatabase(db *gorm.DB) error {
	log.Println("🌱 Seeding database with initial data...")

	// Check if data already exists
	var count int64
	db.Model(&domain.Book{}).Count(&count)
	if count > 0 {
		log.Println("📚 Database already contains data, skipping seed")
		return nil
	}

	// Create some sample publishers
	publishers := []domain.Publisher{
		{
			Name:        "O'Reilly Media",
			Country:     "USA",
			Website:     "https://www.oreilly.com",
			Description: "Technology and business learning platform",
		},
		{
			Name:        "Manning Publications",
			Country:     "USA",
			Website:     "https://www.manning.com",
			Description: "Publisher of computer books for software developers",
		},
		{
			Name:        "Addison-Wesley",
			Country:     "USA",
			Website:     "https://www.awprofessional.com",
			Description: "Publisher of technology and computer science books",
		},
	}

	for i := range publishers {
		if err := db.Create(&publishers[i]).Error; err != nil {
			return fmt.Errorf("failed to create publisher: %w", err)
		}
	}

	// Create some sample authors
	authors := []domain.Author{
		{
			Name:    "Martin Fowler",
			Bio:     "Software development expert and author",
			Country: "UK",
		},
		{
			Name:    "Robert C. Martin",
			Bio:     "Software engineer and instructor, known as Uncle Bob",
			Country: "USA",
		},
		{
			Name:    "Eric Evans",
			Bio:     "Domain-Driven Design pioneer",
			Country: "USA",
		},
		{
			Name:    "Andrew S. Tanenbaum",
			Bio:     "Computer science professor and author",
			Country: "Netherlands",
		},
		{
			Name:    "Maarten van Steen",
			Bio:     "Distributed systems researcher",
			Country: "Netherlands",
		},
	}

	for i := range authors {
		if err := db.Create(&authors[i]).Error; err != nil {
			return fmt.Errorf("failed to create author: %w", err)
		}
	}

	// Create some sample categories
	categories := []domain.Category{
		{
			Name:        "Programming",
			Slug:        "programming",
			Description: "Software development and programming books",
		},
		{
			Name:        "Distributed Systems",
			Slug:        "distributed-systems",
			Description: "Books about distributed systems and architecture",
		},
		{
			Name:        "Software Architecture",
			Slug:        "software-architecture",
			Description: "System design and architecture patterns",
		},
		{
			Name:        "Databases",
			Slug:        "databases",
			Description: "Database design and management",
		},
		{
			Name:        "Cloud Computing",
			Slug:        "cloud-computing",
			Description: "Cloud platforms and services",
		},
	}

	for i := range categories {
		if err := db.Create(&categories[i]).Error; err != nil {
			return fmt.Errorf("failed to create category: %w", err)
		}
	}

	// Create some sample books
	pubDate1 := time.Date(2018, 3, 9, 0, 0, 0, 0, time.UTC)
	pubDate2 := time.Date(2008, 8, 1, 0, 0, 0, 0, time.UTC)
	pubDate3 := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)

	books := []domain.Book{
		{
			ISBN:            "9781492032649",
			Title:           "Building Microservices: Designing Fine-Grained Systems",
			Description:     "A comprehensive guide to building microservices architecture",
			Price:           49.99,
			StockQuantity:   50,
			PublisherID:     &publishers[0].ID,
			CoverImageURL:   "https://covers.oreillystatic.com/images/0636920033158/cat.gif",
			PublicationDate: &pubDate1,
			Language:        "English",
			PageCount:       280,
		},
		{
			ISBN:            "9780132350884",
			Title:           "Clean Code: A Handbook of Agile Software Craftsmanship",
			Description:     "Even bad code can function. But if code isn't clean, it can bring a development organization to its knees.",
			Price:           44.99,
			StockQuantity:   75,
			PublisherID:     &publishers[2].ID,
			CoverImageURL:   "https://m.media-amazon.com/images/I/51E2055ZGUL.jpg",
			PublicationDate: &pubDate2,
			Language:        "English",
			PageCount:       464,
		},
		{
			ISBN:            "9789023456789",
			Title:           "Distributed Systems: Principles and Paradigms",
			Description:     "This book covers the principles, advanced concepts, and technologies of distributed systems in detail.",
			Price:           89.99,
			StockQuantity:   30,
			PublisherID:     &publishers[2].ID,
			CoverImageURL:   "https://m.media-amazon.com/images/I/51OE8WCts7L._SL1360_.jpg",
			PublicationDate: &pubDate3,
			Language:        "English",
			PageCount:       596,
		},
	}

	for i := range books {
		if err := db.Create(&books[i]).Error; err != nil {
			return fmt.Errorf("failed to create book: %w", err)
		}
	}

	// Associate books with authors
	if err := db.Model(&books[0]).Association("Authors").Append(&authors[0]); err != nil {
		return fmt.Errorf("failed to associate book with author: %w", err)
	}

	if err := db.Model(&books[1]).Association("Authors").Append(&authors[1]); err != nil {
		return fmt.Errorf("failed to associate book with author: %w", err)
	}

	if err := db.Model(&books[2]).Association("Authors").Append([]domain.Author{authors[3], authors[4]}); err != nil {
		return fmt.Errorf("failed to associate book with authors: %w", err)
	}

	// Associate books with categories
	if err := db.Model(&books[0]).Association("Categories").Append([]domain.Category{categories[1], categories[2]}); err != nil {
		return fmt.Errorf("failed to associate book with categories: %w", err)
	}

	if err := db.Model(&books[1]).Association("Categories").Append(&categories[0]); err != nil {
		return fmt.Errorf("failed to associate book with category: %w", err)
	}

	if err := db.Model(&books[2]).Association("Categories").Append(&categories[1]); err != nil {
		return fmt.Errorf("failed to associate book with category: %w", err)
	}

	log.Println("Database seeding completed successfully")
	return nil
}
