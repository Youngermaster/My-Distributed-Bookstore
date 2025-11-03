package repository

import (
	"time"

	"gorm.io/gorm"
)

type AdminRepository interface {
	// For future admin-specific database operations
	// Currently, most data comes from other services via gRPC
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{
		db: db,
	}
}

// Helper functions for time range queries
func GetDayStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func GetWeekStart() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
}

func GetMonthStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}
