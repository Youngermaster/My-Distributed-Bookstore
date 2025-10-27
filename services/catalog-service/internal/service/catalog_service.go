package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/domain"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/repository"
)

type CatalogService interface {
	// Books
	CreateBook(ctx context.Context, req CreateBookRequest) (*domain.Book, error)
	GetBook(ctx context.Context, id uuid.UUID) (*domain.Book, error)
	GetBookByISBN(ctx context.Context, isbn string) (*domain.Book, error)
	ListBooks(ctx context.Context, filter repository.BookFilter) (*BookListResponse, error)
	UpdateBook(ctx context.Context, id uuid.UUID, req UpdateBookRequest) (*domain.Book, error)
	DeleteBook(ctx context.Context, id uuid.UUID) error
	SearchBooks(ctx context.Context, query string, filter repository.BookFilter) (*BookListResponse, error)
	UpdateBookStock(ctx context.Context, id uuid.UUID, quantity int) error

	// Authors
	CreateAuthor(ctx context.Context, req CreateAuthorRequest) (*domain.Author, error)
	GetAuthor(ctx context.Context, id uuid.UUID) (*domain.Author, error)
	ListAuthors(ctx context.Context, page, pageSize int) (*AuthorListResponse, error)
	UpdateAuthor(ctx context.Context, id uuid.UUID, req UpdateAuthorRequest) (*domain.Author, error)
	DeleteAuthor(ctx context.Context, id uuid.UUID) error

	// Categories
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*domain.Category, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error)
	ListCategories(ctx context.Context, hierarchical bool) ([]domain.Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req UpdateCategoryRequest) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error

	// Publishers
	CreatePublisher(ctx context.Context, req CreatePublisherRequest) (*domain.Publisher, error)
	GetPublisher(ctx context.Context, id uuid.UUID) (*domain.Publisher, error)
	ListPublishers(ctx context.Context, page, pageSize int) (*PublisherListResponse, error)
	UpdatePublisher(ctx context.Context, id uuid.UUID, req UpdatePublisherRequest) (*domain.Publisher, error)
	DeletePublisher(ctx context.Context, id uuid.UUID) error
}

type catalogService struct {
	bookRepo      repository.BookRepository
	authorRepo    repository.AuthorRepository
	categoryRepo  repository.CategoryRepository
	publisherRepo repository.PublisherRepository
}

func NewCatalogService(
	bookRepo repository.BookRepository,
	authorRepo repository.AuthorRepository,
	categoryRepo repository.CategoryRepository,
	publisherRepo repository.PublisherRepository,
) CatalogService {
	return &catalogService{
		bookRepo:      bookRepo,
		authorRepo:    authorRepo,
		categoryRepo:  categoryRepo,
		publisherRepo: publisherRepo,
	}
}

// Book operations

func (s *catalogService) CreateBook(ctx context.Context, req CreateBookRequest) (*domain.Book, error) {
	// Validate ISBN doesn't already exist
	existing, err := s.bookRepo.GetByISBN(ctx, req.ISBN)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("book with ISBN %s already exists", req.ISBN)
	}

	book := &domain.Book{
		ISBN:            req.ISBN,
		Title:           req.Title,
		Description:     req.Description,
		Price:           req.Price,
		StockQuantity:   req.StockQuantity,
		PublisherID:     req.PublisherID,
		CoverImageURL:   req.CoverImageURL,
		PublicationDate: req.PublicationDate,
		Language:        req.Language,
		PageCount:       req.PageCount,
	}

	// Load authors if provided
	if len(req.AuthorIDs) > 0 {
		for _, authorID := range req.AuthorIDs {
			author, err := s.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return nil, fmt.Errorf("failed to get author: %w", err)
			}
			if author == nil {
				return nil, fmt.Errorf("author with ID %s not found", authorID)
			}
			book.Authors = append(book.Authors, *author)
		}
	}

	// Load categories if provided
	if len(req.CategoryIDs) > 0 {
		for _, categoryID := range req.CategoryIDs {
			category, err := s.categoryRepo.GetByID(ctx, categoryID)
			if err != nil {
				return nil, fmt.Errorf("failed to get category: %w", err)
			}
			if category == nil {
				return nil, fmt.Errorf("category with ID %s not found", categoryID)
			}
			book.Categories = append(book.Categories, *category)
		}
	}

	if err := s.bookRepo.Create(ctx, book); err != nil {
		return nil, err
	}

	log.Printf("✅ Book created: %s (ID: %s)", book.Title, book.ID)
	return book, nil
}

func (s *catalogService) GetBook(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	book, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, fmt.Errorf("book not found")
	}
	return book, nil
}

func (s *catalogService) GetBookByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	book, err := s.bookRepo.GetByISBN(ctx, isbn)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, fmt.Errorf("book not found")
	}
	return book, nil
}

func (s *catalogService) ListBooks(ctx context.Context, filter repository.BookFilter) (*BookListResponse, error) {
	books, total, err := s.bookRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &BookListResponse{
		Books:      books,
		Total:      total,
		Page:       filter.Page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *catalogService) UpdateBook(ctx context.Context, id uuid.UUID, req UpdateBookRequest) (*domain.Book, error) {
	book, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, fmt.Errorf("book not found")
	}

	// Update fields if provided
	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Description != nil {
		book.Description = *req.Description
	}
	if req.Price != nil {
		book.Price = *req.Price
	}
	if req.StockQuantity != nil {
		book.StockQuantity = *req.StockQuantity
	}
	if req.PublisherID != nil {
		book.PublisherID = req.PublisherID
	}
	if req.CoverImageURL != nil {
		book.CoverImageURL = *req.CoverImageURL
	}
	if req.PublicationDate != nil {
		book.PublicationDate = req.PublicationDate
	}
	if req.Language != nil {
		book.Language = *req.Language
	}
	if req.PageCount != nil {
		book.PageCount = *req.PageCount
	}

	if err := s.bookRepo.Update(ctx, book); err != nil {
		return nil, err
	}

	log.Printf("✅ Book updated: %s (ID: %s)", book.Title, book.ID)
	return book, nil
}

func (s *catalogService) DeleteBook(ctx context.Context, id uuid.UUID) error {
	if err := s.bookRepo.Delete(ctx, id); err != nil {
		return err
	}

	log.Printf("🗑️  Book deleted: ID %s", id)
	return nil
}

func (s *catalogService) SearchBooks(ctx context.Context, query string, filter repository.BookFilter) (*BookListResponse, error) {
	books, total, err := s.bookRepo.Search(ctx, query, filter)
	if err != nil {
		return nil, err
	}

	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &BookListResponse{
		Books:      books,
		Total:      total,
		Page:       filter.Page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *catalogService) UpdateBookStock(ctx context.Context, id uuid.UUID, quantity int) error {
	if err := s.bookRepo.UpdateStock(ctx, id, quantity); err != nil {
		return err
	}

	log.Printf("📦 Book stock updated: ID %s, new quantity: %d", id, quantity)
	return nil
}

// Author operations

func (s *catalogService) CreateAuthor(ctx context.Context, req CreateAuthorRequest) (*domain.Author, error) {
	author := &domain.Author{
		Name:      req.Name,
		Bio:       req.Bio,
		BirthDate: req.BirthDate,
		Country:   req.Country,
		ImageURL:  req.ImageURL,
	}

	if err := s.authorRepo.Create(ctx, author); err != nil {
		return nil, err
	}

	log.Printf("✅ Author created: %s (ID: %s)", author.Name, author.ID)
	return author, nil
}

func (s *catalogService) GetAuthor(ctx context.Context, id uuid.UUID) (*domain.Author, error) {
	author, err := s.authorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, fmt.Errorf("author not found")
	}
	return author, nil
}

func (s *catalogService) ListAuthors(ctx context.Context, page, pageSize int) (*AuthorListResponse, error) {
	authors, total, err := s.authorRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	if pageSize < 1 {
		pageSize = 20
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &AuthorListResponse{
		Authors:    authors,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *catalogService) UpdateAuthor(ctx context.Context, id uuid.UUID, req UpdateAuthorRequest) (*domain.Author, error) {
	author, err := s.authorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, fmt.Errorf("author not found")
	}

	if req.Name != nil {
		author.Name = *req.Name
	}
	if req.Bio != nil {
		author.Bio = *req.Bio
	}
	if req.BirthDate != nil {
		author.BirthDate = req.BirthDate
	}
	if req.Country != nil {
		author.Country = *req.Country
	}
	if req.ImageURL != nil {
		author.ImageURL = *req.ImageURL
	}

	if err := s.authorRepo.Update(ctx, author); err != nil {
		return nil, err
	}

	log.Printf("✅ Author updated: %s (ID: %s)", author.Name, author.ID)
	return author, nil
}

func (s *catalogService) DeleteAuthor(ctx context.Context, id uuid.UUID) error {
	if err := s.authorRepo.Delete(ctx, id); err != nil {
		return err
	}

	log.Printf("🗑️  Author deleted: ID %s", id)
	return nil
}

// Category operations

func (s *catalogService) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*domain.Category, error) {
	category := &domain.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	log.Printf("✅ Category created: %s (ID: %s)", category.Name, category.ID)
	return category, nil
}

func (s *catalogService) GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}

func (s *catalogService) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}

func (s *catalogService) ListCategories(ctx context.Context, hierarchical bool) ([]domain.Category, error) {
	if hierarchical {
		return s.categoryRepo.ListHierarchical(ctx)
	}
	return s.categoryRepo.List(ctx)
}

func (s *catalogService) UpdateCategory(ctx context.Context, id uuid.UUID, req UpdateCategoryRequest) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	log.Printf("✅ Category updated: %s (ID: %s)", category.Name, category.ID)
	return category, nil
}

func (s *catalogService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return err
	}

	log.Printf("🗑️  Category deleted: ID %s", id)
	return nil
}

// Publisher operations

func (s *catalogService) CreatePublisher(ctx context.Context, req CreatePublisherRequest) (*domain.Publisher, error) {
	publisher := &domain.Publisher{
		Name:        req.Name,
		Country:     req.Country,
		Website:     req.Website,
		Description: req.Description,
	}

	if err := s.publisherRepo.Create(ctx, publisher); err != nil {
		return nil, err
	}

	log.Printf("✅ Publisher created: %s (ID: %s)", publisher.Name, publisher.ID)
	return publisher, nil
}

func (s *catalogService) GetPublisher(ctx context.Context, id uuid.UUID) (*domain.Publisher, error) {
	publisher, err := s.publisherRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if publisher == nil {
		return nil, fmt.Errorf("publisher not found")
	}
	return publisher, nil
}

func (s *catalogService) ListPublishers(ctx context.Context, page, pageSize int) (*PublisherListResponse, error) {
	publishers, total, err := s.publisherRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	if pageSize < 1 {
		pageSize = 20
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PublisherListResponse{
		Publishers: publishers,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *catalogService) UpdatePublisher(ctx context.Context, id uuid.UUID, req UpdatePublisherRequest) (*domain.Publisher, error) {
	publisher, err := s.publisherRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if publisher == nil {
		return nil, fmt.Errorf("publisher not found")
	}

	if req.Name != nil {
		publisher.Name = *req.Name
	}
	if req.Country != nil {
		publisher.Country = *req.Country
	}
	if req.Website != nil {
		publisher.Website = *req.Website
	}
	if req.Description != nil {
		publisher.Description = *req.Description
	}

	if err := s.publisherRepo.Update(ctx, publisher); err != nil {
		return nil, err
	}

	log.Printf("✅ Publisher updated: %s (ID: %s)", publisher.Name, publisher.ID)
	return publisher, nil
}

func (s *catalogService) DeletePublisher(ctx context.Context, id uuid.UUID) error {
	if err := s.publisherRepo.Delete(ctx, id); err != nil {
		return err
	}

	log.Printf("🗑️  Publisher deleted: ID %s", id)
	return nil
}
