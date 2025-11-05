package http

import (
	"github.com/google/uuid"
)

// parseUUID parses a string to UUID
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
