package notification

import (
	"fmt"
	"strings"
	"time"

	"github.com/youngermaster/distributed-bookstore/notification-service/internal/events"
)

type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelPush  Channel = "push"
)

type Message struct {
	ID        string
	EventType string
	Channel   Channel
	Recipient string
	Subject   string
	Body      string
	SentAt    time.Time
	Metadata  map[string]interface{}
}

func (m Message) String() string {
	return fmt.Sprintf("[%s] %s -> %s", m.Channel, m.EventType, m.Recipient)
}

func BuildMessages(evt events.NotificationEvent) []Message {
	switch evt.EventType {
	case "user.registered":
		return []Message{
			{
				ID:        evt.EventID,
				EventType: evt.EventType,
				Channel:   ChannelEmail,
				Recipient: firstNonEmpty(evt.Email, extractString(evt.Payload, "email")),
				Subject:   "Welcome to Ohara Bookstore!",
				Body: fmt.Sprintf("Hi %s,\n\nThanks for joining Ohara Bookstore. We're excited to help you discover your next favorite read.\n\nHappy reading!\n— The Ohara Team",
					extractString(evt.Payload, "full_name")),
				SentAt:   time.Now(),
				Metadata: evt.Payload,
			},
		}
	case "user.logged_in":
		return []Message{
			{
				ID:        evt.EventID,
				EventType: evt.EventType,
				Channel:   ChannelEmail,
				Recipient: firstNonEmpty(evt.Email, extractString(evt.Payload, "email")),
				Subject:   "New login detected",
				Body: fmt.Sprintf("Hello %s,\n\nWe noticed a new login to your account at %s. If this was you, you can safely ignore this message. Otherwise, reset your password immediately.",
					extractString(evt.Payload, "full_name"),
					time.Now().Format(time.RFC1123)),
				SentAt:   time.Now(),
				Metadata: evt.Payload,
			},
		}
	case "wishlist.added":
		return []Message{
			{
				ID:        evt.EventID,
				EventType: evt.EventType,
				Channel:   ChannelPush,
				Recipient: extractString(evt.Payload, "user_id"),
				Subject:   "Wishlist update",
				Body: fmt.Sprintf("“%s” is now in your wishlist. We'll let you know about price drops!",
					extractString(evt.Payload, "book_title")),
				SentAt:   time.Now(),
				Metadata: evt.Payload,
			},
		}
	case "cart.item_added":
		return []Message{
			{
				ID:        evt.EventID,
				EventType: evt.EventType,
				Channel:   ChannelPush,
				Recipient: extractString(evt.Payload, "cart_id"),
				Subject:   "Cart updated",
				Body: fmt.Sprintf("Added %s to your cart. You now have %d item(s) ready to checkout.",
					extractString(evt.Payload, "book_title"),
					intFromPayload(evt.Payload, "total_items")),
				SentAt:   time.Now(),
				Metadata: evt.Payload,
			},
		}
	default:
		return []Message{
			{
				ID:        evt.EventID,
				EventType: evt.EventType,
				Channel:   ChannelPush,
				Recipient: extractString(evt.Payload, "user_id"),
				Subject:   fmt.Sprintf("Event: %s", evt.EventType),
				Body:      "A notification event occurred.",
				SentAt:    time.Now(),
				Metadata:  evt.Payload,
			},
		}
	}
}

func extractString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key]; ok {
		switch value := v.(type) {
		case string:
			return value
		case fmt.Stringer:
			return value.String()
		default:
			return ""
		}
	}
	return ""
}

func intFromPayload(m map[string]interface{}, key string) int {
	if m == nil {
		return 0
	}
	if v, ok := m[key]; ok {
		switch value := v.(type) {
		case float64:
			return int(value)
		case int:
			return value
		case int64:
			return int(value)
		}
	}
	return 0
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
