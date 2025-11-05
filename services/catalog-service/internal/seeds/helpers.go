package seeds

import (
	"time"

	"github.com/google/uuid"
)

// UUIDFromString returns a stable UUID generated from the provided value.
// Using a namespace-based UUID keeps IDs consistent across services so that
// relationships (like recommendation data) can reference the same records.
func UUIDFromString(value string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(value))
}

// MustParseDate parses an ISO-8601 date (yyyy-mm-dd) and panics if invalid.
// Seed data is static, so it is better to fail fast during application startup.
func MustParseDate(value string) time.Time {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		panic(err)
	}
	return t
}

// TimePtr returns a pointer to the provided time value.
func TimePtr(t time.Time) *time.Time {
	return &t
}
