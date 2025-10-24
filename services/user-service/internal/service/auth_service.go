package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/domain"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/dto"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/repository"
	"github.com/youngermaster/distributed-bookstore/user-service/pkg/jwt"
	"github.com/youngermaster/distributed-bookstore/user-service/pkg/password"
)

var (
	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo     *repository.UserRepository
	roleRepo     *repository.RoleRepository
	sessionRepo  *repository.SessionRepository
	jwtService   *jwt.JWTService
	passwordSvc  *password.Service
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	sessionRepo *repository.SessionRepository,
	jwtService *jwt.JWTService,
	passwordSvc *password.Service,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		sessionRepo: sessionRepo,
		jwtService:  jwtService,
		passwordSvc: passwordSvc,
	}
}

// Register registers a new user
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Hash password
	passwordHash, err := s.passwordSvc.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		Phone:        req.Phone,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Assign default "customer" role
	customerRole, err := s.roleRepo.GetByName("customer")
	if err == nil {
		s.userRepo.AssignRole(user.ID, customerRole.ID)
	}

	// Reload user with roles
	user, err = s.userRepo.GetByID(user.ID)
	if err != nil {
		return nil, err
	}

	return s.mapUserToResponse(user), nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if err := s.passwordSvc.Verify(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Extract role names
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	// Generate tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, roleNames)
	if err != nil {
		return nil, err
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	return &dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		TokenType:    tokenPair.TokenType,
		User:         *s.mapUserToResponse(user),
	}, nil
}

// RefreshToken generates new tokens from a refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*dto.LoginResponse, error) {
	// Validate refresh token
	userID, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Extract role names
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	// Generate new tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, roleNames)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		TokenType:    tokenPair.TokenType,
		User:         *s.mapUserToResponse(user),
	}, nil
}

// ValidateToken validates an access token and returns claims
func (s *AuthService) ValidateToken(token string) (*jwt.Claims, error) {
	return s.jwtService.ValidateToken(token)
}

// Logout logs out a user (invalidates sessions)
func (s *AuthService) Logout(userID uuid.UUID) error {
	return s.sessionRepo.DeleteByUserID(userID)
}

// mapUserToResponse converts domain user to response DTO
func (s *AuthService) mapUserToResponse(user *domain.User) *dto.UserResponse {
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
