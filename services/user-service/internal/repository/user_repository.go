package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/domain"
	"gorm.io/gorm"
)

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailExists is returned when email already exists
	ErrEmailExists = errors.New("email already exists")
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *domain.User) error {
	// Check if email already exists
	var count int64
	if err := r.db.Model(&domain.User{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrEmailExists
	}

	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Roles").Preload("Addresses").First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&domain.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// AssignRole assigns a role to a user
func (r *UserRepository) AssignRole(userID, roleID uuid.UUID) error {
	var user domain.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	var role domain.Role
	if err := r.db.First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&user).Association("Roles").Append(&role)
}

// RemoveRole removes a role from a user
func (r *UserRepository) RemoveRole(userID, roleID uuid.UUID) error {
	var user domain.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	var role domain.Role
	if err := r.db.First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&user).Association("Roles").Delete(&role)
}

// GetUserRoles retrieves all roles for a user
func (r *UserRepository) GetUserRoles(userID uuid.UUID) ([]domain.Role, error) {
	var user domain.User
	if err := r.db.Preload("Roles").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return user.Roles, nil
}

// AddressRepository handles database operations for addresses
type AddressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new address repository
func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

// Create creates a new address
func (r *AddressRepository) Create(address *domain.Address) error {
	// If this is set as default, unset other defaults for this user
	if address.IsDefault {
		r.db.Model(&domain.Address{}).Where("user_id = ? AND is_default = ?", address.UserID, true).
			Update("is_default", false)
	}
	return r.db.Create(address).Error
}

// GetByID retrieves an address by ID
func (r *AddressRepository) GetByID(id uuid.UUID) (*domain.Address, error) {
	var address domain.Address
	err := r.db.First(&address, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		return nil, err
	}
	return &address, nil
}

// GetByUserID retrieves all addresses for a user
func (r *AddressRepository) GetByUserID(userID uuid.UUID) ([]domain.Address, error) {
	var addresses []domain.Address
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error
	return addresses, err
}

// Update updates an address
func (r *AddressRepository) Update(address *domain.Address) error {
	// If this is set as default, unset other defaults for this user
	if address.IsDefault {
		r.db.Model(&domain.Address{}).
			Where("user_id = ? AND id != ? AND is_default = ?", address.UserID, address.ID, true).
			Update("is_default", false)
	}
	return r.db.Save(address).Error
}

// Delete soft deletes an address
func (r *AddressRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&domain.Address{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("address not found")
	}
	return nil
}

// SessionRepository handles database operations for sessions
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session
func (r *SessionRepository) Create(session *domain.Session) error {
	return r.db.Create(session).Error
}

// GetByTokenHash retrieves a session by token hash
func (r *SessionRepository) GetByTokenHash(tokenHash string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.Where("token_hash = ?", tokenHash).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	return &session, nil
}

// DeleteByUserID deletes all sessions for a user (logout all devices)
func (r *SessionRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.Session{}).Error
}

// DeleteExpired deletes all expired sessions
func (r *SessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < NOW()").Delete(&domain.Session{}).Error
}

// RoleRepository handles database operations for roles
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create creates a new role
func (r *RoleRepository) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}

// GetByID retrieves a role by ID
func (r *RoleRepository) GetByID(id uuid.UUID) (*domain.Role, error) {
	var role domain.Role
	err := r.db.First(&role, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

// GetByName retrieves a role by name
func (r *RoleRepository) GetByName(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

// GetAll retrieves all roles
func (r *RoleRepository) GetAll() ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.Find(&roles).Error
	return roles, err
}
