package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/domain"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/dto"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/repository"
	"github.com/youngermaster/distributed-bookstore/user-service/pkg/password"
)

// UserService handles user management business logic
type UserService struct {
	userRepo    *repository.UserRepository
	addressRepo *repository.AddressRepository
	passwordSvc *password.Service
}

// NewUserService creates a new user service
func NewUserService(
	userRepo *repository.UserRepository,
	addressRepo *repository.AddressRepository,
	passwordSvc *password.Service,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		addressRepo: addressRepo,
		passwordSvc: passwordSvc,
	}
}

// GetProfile gets user profile by ID
func (s *UserService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return mapUserToResponse(user), nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return mapUserToResponse(user), nil
}

// ChangePassword changes user password
func (s *UserService) ChangePassword(userID uuid.UUID, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := s.passwordSvc.Verify(req.CurrentPassword, user.PasswordHash); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	newHash, err := s.passwordSvc.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash
	return s.userRepo.Update(user)
}

// GetAddresses retrieves all addresses for a user
func (s *UserService) GetAddresses(userID uuid.UUID) ([]dto.AddressResponse, error) {
	addresses, err := s.addressRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.AddressResponse, len(addresses))
	for i, addr := range addresses {
		response[i] = mapAddressToResponse(&addr)
	}

	return response, nil
}

// CreateAddress creates a new address
func (s *UserService) CreateAddress(userID uuid.UUID, req dto.AddressRequest) (*dto.AddressResponse, error) {
	address := &domain.Address{
		UserID:       userID,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		Country:      req.Country,
		IsDefault:    req.IsDefault,
	}

	if err := s.addressRepo.Create(address); err != nil {
		return nil, err
	}

	return &dto.AddressResponse{
		ID:           address.ID,
		UserID:       address.UserID,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		State:        address.State,
		PostalCode:   address.PostalCode,
		Country:      address.Country,
		IsDefault:    address.IsDefault,
		CreatedAt:    address.CreatedAt,
		UpdatedAt:    address.UpdatedAt,
	}, nil
}

// UpdateAddress updates an existing address
func (s *UserService) UpdateAddress(userID, addressID uuid.UUID, req dto.AddressRequest) (*dto.AddressResponse, error) {
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if address.UserID != userID {
		return nil, errors.New("address not found")
	}

	// Update fields
	address.AddressLine1 = req.AddressLine1
	address.AddressLine2 = req.AddressLine2
	address.City = req.City
	address.State = req.State
	address.PostalCode = req.PostalCode
	address.Country = req.Country
	address.IsDefault = req.IsDefault

	if err := s.addressRepo.Update(address); err != nil {
		return nil, err
	}

	return &dto.AddressResponse{
		ID:           address.ID,
		UserID:       address.UserID,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		State:        address.State,
		PostalCode:   address.PostalCode,
		Country:      address.Country,
		IsDefault:    address.IsDefault,
		CreatedAt:    address.CreatedAt,
		UpdatedAt:    address.UpdatedAt,
	}, nil
}

// DeleteAddress deletes an address
func (s *UserService) DeleteAddress(userID, addressID uuid.UUID) error {
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return err
	}

	// Verify ownership
	if address.UserID != userID {
		return errors.New("address not found")
	}

	return s.addressRepo.Delete(addressID)
}

// Helper functions

func mapUserToResponse(user *domain.User) *dto.UserResponse {
	roles := make([]dto.RoleResponse, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		}
	}

	return &dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		Phone:       user.Phone,
		Roles:       roles,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}
}

func mapAddressToResponse(address *domain.Address) dto.AddressResponse {
	return dto.AddressResponse{
		ID:           address.ID,
		UserID:       address.UserID,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		State:        address.State,
		PostalCode:   address.PostalCode,
		Country:      address.Country,
		IsDefault:    address.IsDefault,
		CreatedAt:    address.CreatedAt,
		UpdatedAt:    address.UpdatedAt,
	}
}
