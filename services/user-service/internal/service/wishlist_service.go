package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/domain"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/dto"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/events"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/repository"
)

// WishlistService handles wishlist business logic
type WishlistService struct {
	wishlistRepo *repository.WishlistRepository
	eventPub     *events.Publisher
	// TODO: Add catalog service client for fetching book details
	// catalogClient *catalog.Client
}

// NewWishlistService creates a new wishlist service
func NewWishlistService(wishlistRepo *repository.WishlistRepository, eventPublisher *events.Publisher) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
		eventPub:     eventPublisher,
	}
}

// Add adds a book to user's wishlist
func (s *WishlistService) Add(userID, bookID uuid.UUID) (*dto.WishlistResponse, error) {
	// TODO: Validate that book exists by calling catalog service

	wishlist, err := s.wishlistRepo.Add(userID, bookID)
	if err != nil {
		return nil, err
	}

	response := mapWishlistToResponse(wishlist)

	s.publishEvent(events.NewNotificationEvent("wishlist.added", "user-service").
		WithUser(userID, "").
		WithPayload(map[string]interface{}{
			"wishlist_id": wishlist.ID.String(),
			"book_id":     bookID.String(),
		}))

	return response, nil
}

// Remove removes a book from user's wishlist
func (s *WishlistService) Remove(userID, bookID uuid.UUID) error {
	if err := s.wishlistRepo.Remove(userID, bookID); err != nil {
		return err
	}

	s.publishEvent(events.NewNotificationEvent("wishlist.removed", "user-service").
		WithUser(userID, "").
		WithPayload(map[string]interface{}{
			"book_id": bookID.String(),
		}))

	return nil
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
	if err := s.wishlistRepo.Clear(userID); err != nil {
		return err
	}

	s.publishEvent(events.NewNotificationEvent("wishlist.cleared", "user-service").
		WithUser(userID, ""))

	return nil
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

func (s *WishlistService) publishEvent(evt events.NotificationEvent) {
	if s.eventPub == nil {
		return
	}

	if err := s.eventPub.Publish(context.Background(), evt.EventType, evt); err != nil {
		log.Printf("failed to publish wishlist event %s: %v", evt.EventType, err)
	}
}

// TODO: Implement when gRPC client is ready
// func (s *WishlistService) fetchBookDetails(bookID uuid.UUID) *dto.BookSummary {
// 	// Call catalog service via gRPC to get book details
// 	return nil
// }
