package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents an active user session
type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string    `gorm:"size:255;not null;index" json:"-"`
	IPAddress string    `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent string    `gorm:"type:text" json:"user_agent,omitempty"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

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
