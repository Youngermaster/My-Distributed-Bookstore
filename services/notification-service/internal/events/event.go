package events

import (
	"encoding/json"
	"time"
)

type NotificationEvent struct {
	EventID    string                 `json:"event_id"`
	EventType  string                 `json:"event_type"`
	Source     string                 `json:"source"`
	OccurredAt time.Time              `json:"occurred_at"`
	UserID     string                 `json:"user_id,omitempty"`
	Email      string                 `json:"email,omitempty"`
	Phone      string                 `json:"phone,omitempty"`
	Channel    string                 `json:"channel,omitempty"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
}

func FromJSON(data []byte) (NotificationEvent, error) {
	var evt NotificationEvent
	err := json.Unmarshal(data, &evt)
	return evt, err
}

func (e NotificationEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
