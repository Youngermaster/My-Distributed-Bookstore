package dto

import (
	"time"

	"github.com/google/uuid"
)

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a login response with tokens
type LoginResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresAt    time.Time  `json:"expires_at"`
	TokenType    string     `json:"token_type"`
	User         UserResponse `json:"user"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FullName string `json:"full_name,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID          uuid.UUID       `json:"id"`
	Email       string          `json:"email"`
	FullName    string          `json:"full_name"`
	Phone       string          `json:"phone,omitempty"`
	Roles       []RoleResponse  `json:"roles,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	LastLoginAt *time.Time      `json:"last_login_at,omitempty"`
}

// RoleResponse represents a role in API responses
type RoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
}

// AddressRequest represents a request to create/update an address
type AddressRequest struct {
	AddressLine1 string `json:"address_line1" validate:"required"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city" validate:"required"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postal_code" validate:"required"`
	Country      string `json:"country" validate:"required"`
	IsDefault    bool   `json:"is_default"`
}

// AddressResponse represents an address in API responses
type AddressResponse struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	AddressLine1 string    `json:"address_line1"`
	AddressLine2 string    `json:"address_line2,omitempty"`
	City         string    `json:"city"`
	State        string    `json:"state,omitempty"`
	PostalCode   string    `json:"postal_code"`
	Country      string    `json:"country"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
