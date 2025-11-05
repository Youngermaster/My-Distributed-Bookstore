package service

import (
	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/domain"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/dto"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/repository"
)

// WishlistService handles wishlist business logic
type WishlistService struct {
	wishlistRepo *repository.WishlistRepository
	// TODO: Add catalog service client for fetching book details
	// catalogClient *catalog.Client
}

// NewWishlistService creates a new wishlist service
func NewWishlistService(wishlistRepo *repository.WishlistRepository) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
	}
}

// Add adds a book to user's wishlist
func (s *WishlistService) Add(userID, bookID uuid.UUID) (*dto.WishlistResponse, error) {
	// TODO: Validate that book exists by calling catalog service

	wishlist, err := s.wishlistRepo.Add(userID, bookID)
	if err != nil {
		return nil, err
	}

	return mapWishlistToResponse(wishlist), nil
}

// Remove removes a book from user's wishlist
func (s *WishlistService) Remove(userID, bookID uuid.UUID) error {
	return s.wishlistRepo.Remove(userID, bookID)
}

// GetAll retrieves all wishlist items for a user
func (s *WishlistService) GetAll(userID uuid.UUID) ([]dto.WishlistResponse, error) {
	wishlists, err := s.wishlistRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.WishlistResponse, len(wishlists))
	for i, wishlist := range wishlists {
		response[i] = *mapWishlistToResponse(&wishlist)
	}

	return response, nil
}

// GetBookIDs retrieves all book IDs in user's wishlist
func (s *WishlistService) GetBookIDs(userID uuid.UUID) ([]uuid.UUID, error) {
	return s.wishlistRepo.GetBookIDs(userID)
}

// GetAllWithBooks retrieves all wishlist items with book details
func (s *WishlistService) GetAllWithBooks(userID uuid.UUID) ([]dto.WishlistWithBookResponse, error) {
	wishlists, err := s.wishlistRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.WishlistWithBookResponse, len(wishlists))
	for i, wishlist := range wishlists {
		response[i] = dto.WishlistWithBookResponse{
			ID:        wishlist.ID,
			UserID:    wishlist.UserID,
			BookID:    wishlist.BookID,
			CreatedAt: wishlist.CreatedAt,
			// TODO: Fetch book details from catalog service
			// Book: s.fetchBookDetails(wishlist.BookID),
		}
	}

	return response, nil
}

// Exists checks if a book is in user's wishlist
func (s *WishlistService) Exists(userID, bookID uuid.UUID) (bool, error) {
	return s.wishlistRepo.Exists(userID, bookID)
}

// Clear removes all items from user's wishlist
func (s *WishlistService) Clear(userID uuid.UUID) error {
	return s.wishlistRepo.Clear(userID)
}

// Helper functions

func mapWishlistToResponse(wishlist *domain.Wishlist) *dto.WishlistResponse {
	return &dto.WishlistResponse{
		ID:        wishlist.ID,
		UserID:    wishlist.UserID,
		BookID:    wishlist.BookID,
		CreatedAt: wishlist.CreatedAt,
	}
}

// TODO: Implement when gRPC client is ready
// func (s *WishlistService) fetchBookDetails(bookID uuid.UUID) *dto.BookSummary {
// 	// Call catalog service via gRPC to get book details
// 	return nil
// }
