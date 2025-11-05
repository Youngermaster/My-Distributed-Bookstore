package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NotificationEvent struct {
	EventID    string                 `json:"event_id"`
	EventType  string                 `json:"event_type"`
	Source     string                 `json:"source"`
	OccurredAt time.Time              `json:"occurred_at"`
	UserID     string                 `json:"user_id,omitempty"`
	Email      string                 `json:"email,omitempty"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
}

func NewNotificationEvent(eventType, source string) NotificationEvent {
	return NotificationEvent{
		EventID:    uuid.NewString(),
		EventType:  eventType,
		Source:     source,
		OccurredAt: time.Now().UTC(),
		Payload:    make(map[string]interface{}),
	}
}

func (e NotificationEvent) WithPayload(data map[string]interface{}) NotificationEvent {
	if e.Payload == nil {
		e.Payload = make(map[string]interface{})
	}
	for k, v := range data {
		e.Payload[k] = v
	}
	return e
}

func (e NotificationEvent) Marshal() ([]byte, error) {
	return json.Marshal(e)
}
