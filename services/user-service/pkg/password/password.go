package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidPassword is returned when password verification fails
	ErrInvalidPassword = errors.New("invalid password")
	// ErrPasswordTooShort is returned when password is too short
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)

const (
	// MinPasswordLength is the minimum password length
	MinPasswordLength = 8
	// DefaultCost is the default bcrypt cost
	DefaultCost = 12
)

// Service handles password hashing and verification
type Service struct {
	cost int
}

// NewService creates a new password service
func NewService(cost int) *Service {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = DefaultCost
	}
	return &Service{
		cost: cost,
	}
}

// Hash creates a bcrypt hash from a plain text password
func (s *Service) Hash(password string) (string, error) {
	if len(password) < MinPasswordLength {
		return "", ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Verify compares a plain text password with a hash
func (s *Service) Verify(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}
		return err
	}
	return nil
}

// NeedsRehash checks if a password hash needs to be rehashed
func (s *Service) NeedsRehash(hash string) bool {
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return false
	}
	return cost != s.cost
}

// ValidatePasswordStrength validates password meets minimum requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	// Add additional validation rules here:
	// - Check for at least one uppercase letter
	// - Check for at least one lowercase letter
	// - Check for at least one digit
	// - Check for at least one special character

	return nil
}
