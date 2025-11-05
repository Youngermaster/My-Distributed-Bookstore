package config

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/domain"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/seeds"
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
	if err := db.Model(&domain.Book{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing data: %w", err)
	}
	if count > 0 {
		log.Println("📚 Database already contains data, skipping seed")
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Seed categories
		categorySeeds := seeds.GetCategories()
		categories := make([]domain.Category, len(categorySeeds))
		for i, seed := range categorySeeds {
			categories[i] = domain.Category{
				ID:          seed.ID,
				Name:        seed.Name,
				Slug:        seed.Slug,
				Description: seed.Description,
			}
		}
		if len(categories) > 0 {
			if err := tx.Create(&categories).Error; err != nil {
				return fmt.Errorf("failed to create categories: %w", err)
			}
		}

		categoryMap := make(map[string]*domain.Category, len(categories))
		for i := range categories {
			category := &categories[i]
			categoryMap[category.Slug] = category
		}

		// Seed authors
		authorSeeds := seeds.GetAuthors()
		authors := make([]domain.Author, len(authorSeeds))
		for i, seed := range authorSeeds {
			authors[i] = domain.Author{
				ID:       seed.ID,
				Name:     seed.Name,
				Bio:      seed.Bio,
				Country:  seed.Country,
				ImageURL: seed.ImageURL,
			}
		}
		if len(authors) > 0 {
			if err := tx.Create(&authors).Error; err != nil {
				return fmt.Errorf("failed to create authors: %w", err)
			}
		}

		authorMap := make(map[string]*domain.Author, len(authors))
		for i, seed := range authorSeeds {
			authorMap[seed.Code] = &authors[i]
		}

		// Seed publishers
		publisherSeeds := seeds.GetPublishers()
		publishers := make([]domain.Publisher, len(publisherSeeds))
		for i, seed := range publisherSeeds {
			publishers[i] = domain.Publisher{
				ID:          seed.ID,
				Name:        seed.Name,
				Country:     seed.Country,
				Website:     seed.Website,
				Description: seed.Description,
			}
		}
		if len(publishers) > 0 {
			if err := tx.Create(&publishers).Error; err != nil {
				return fmt.Errorf("failed to create publishers: %w", err)
			}
		}

		publisherMap := make(map[string]*domain.Publisher, len(publishers))
		for i, seed := range publisherSeeds {
			publisherMap[seed.Code] = &publishers[i]
		}

		// Seed books and relationships
		bookSeeds := seeds.GetBooks()
		books := make([]domain.Book, len(bookSeeds))
		for i, seed := range bookSeeds {
			pubDate := seeds.MustParseDate(seed.PublicationDate)
			coverURL := seed.CoverImageURL
			if coverURL == "" && seed.ISBN != "" {
				coverURL = fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-L.jpg", seed.ISBN)
			}

			books[i] = domain.Book{
				ID:              seed.ID,
				ISBN:            seed.ISBN,
				Title:           seed.Title,
				Description:     seed.Description,
				Price:           seed.Price,
				StockQuantity:   seed.Stock,
				CoverImageURL:   coverURL,
				PublicationDate: seeds.TimePtr(pubDate),
				Language:        seed.Language,
				PageCount:       seed.PageCount,
			}

			if publisher, ok := publisherMap[seed.PublisherCode]; ok {
				books[i].PublisherID = &publisher.ID
			} else {
				return fmt.Errorf("publisher code %q not found for book %q", seed.PublisherCode, seed.Title)
			}
		}

		if len(books) > 0 {
			if err := tx.Create(&books).Error; err != nil {
				return fmt.Errorf("failed to create books: %w", err)
			}
		}

		for i, seed := range bookSeeds {
			book := &books[i]

			// Attach categories
			var categoryRefs []*domain.Category
			seenCategories := make(map[uuid.UUID]struct{})
			for _, slug := range seed.CategorySlugs {
				category, ok := categoryMap[slug]
				if !ok {
					return fmt.Errorf("book %q references unknown category %q", seed.Title, slug)
				}
				if _, exists := seenCategories[category.ID]; !exists {
					seenCategories[category.ID] = struct{}{}
					categoryRefs = append(categoryRefs, category)
				}
			}

			if len(categoryRefs) > 0 {
				assoc := tx.Model(book).Association("Categories")
				for _, category := range categoryRefs {
					if err := assoc.Append(category); err != nil {
						return fmt.Errorf("failed to associate categories for %q: %w", seed.Title, err)
					}
				}
			}

			// Attach authors
			var authorRefs []*domain.Author
			seenAuthors := make(map[uuid.UUID]struct{})
			for _, code := range seed.AuthorCodes {
				author, ok := authorMap[code]
				if !ok {
					return fmt.Errorf("book %q references unknown author %q", seed.Title, code)
				}
				if _, exists := seenAuthors[author.ID]; !exists {
					seenAuthors[author.ID] = struct{}{}
					authorRefs = append(authorRefs, author)
				}
			}

			if len(authorRefs) > 0 {
				assoc := tx.Model(book).Association("Authors")
				for _, author := range authorRefs {
					if err := assoc.Append(author); err != nil {
						return fmt.Errorf("failed to associate authors for %q: %w", seed.Title, err)
					}
				}
			}
		}

		return nil
	})
}
