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
