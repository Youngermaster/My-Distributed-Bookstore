package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user account in the system
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	FullName     string         `gorm:"size:255" json:"full_name"`
	Phone        string         `gorm:"size:20" json:"phone,omitempty"`
	CreatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Roles     []Role    `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	Addresses []Address `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
}

// TableName overrides the table name for GORM
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Role represents a user role for RBAC
type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null;size:50" json:"name"`
	Description string    `gorm:"size:500" json:"description,omitempty"`
	Permissions string    `gorm:"type:jsonb" json:"permissions"` // JSONB array of permission strings
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	Users []User `gorm:"many2many:user_roles;" json:"-"`
}

// TableName overrides the table name for GORM
func (Role) TableName() string {
	return "roles"
}

// BeforeCreate hook to generate UUID
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// Address represents a shipping or billing address
type Address struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	AddressLine1 string         `gorm:"size:255;not null" json:"address_line1"`
	AddressLine2 string         `gorm:"size:255" json:"address_line2,omitempty"`
	City         string         `gorm:"size:100;not null" json:"city"`
	State        string         `gorm:"size:100" json:"state,omitempty"`
	PostalCode   string         `gorm:"size:20;not null" json:"postal_code"`
	Country      string         `gorm:"size:100;not null" json:"country"`
	IsDefault    bool           `gorm:"default:false" json:"is_default"`
	CreatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the table name for GORM
func (Address) TableName() string {
	return "addresses"
}

// BeforeCreate hook to generate UUID
func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// Session represents an active user session
type Session struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string     `gorm:"size:255;not null;index" json:"-"`
	IPAddress string     `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent string     `gorm:"type:text" json:"user_agent,omitempty"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the table name for GORM
func (Session) TableName() string {
	return "sessions"
}

// BeforeCreate hook to generate UUID
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// Wishlist represents a user's wishlist item (bookmarked book)
type Wishlist struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_wishlist_user_book,priority:1" json:"user_id"`
	BookID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_wishlist_user_book,priority:2" json:"book_id"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the table name for GORM
func (Wishlist) TableName() string {
	return "wishlists"
}

// BeforeCreate hook to generate UUID
func (w *Wishlist) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
