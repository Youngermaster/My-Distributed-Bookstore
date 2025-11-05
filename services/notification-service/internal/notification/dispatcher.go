package notification

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/youngermaster/distributed-bookstore/notification-service/internal/events"
)

type Dispatcher struct {
	logger  zerolog.Logger
	senders map[Channel]Sender
	history *History
}

func NewDispatcher(logger zerolog.Logger, senders []Sender, history *History) *Dispatcher {
	senderMap := make(map[Channel]Sender, len(senders))
	for _, sender := range senders {
		senderMap[sender.Channel()] = sender
	}

	return &Dispatcher{
		logger:  logger,
		senders: senderMap,
		history: history,
	}
}

func (d *Dispatcher) HandleEvent(ctx context.Context, evt events.NotificationEvent) error {
	messages := BuildMessages(evt)
	if len(messages) == 0 {
		d.logger.Warn().
			Str("event_type", evt.EventType).
			Msg("no notifications generated for event")
		return nil
	}

	for _, msg := range messages {
		sender, ok := d.senders[msg.Channel]
		if !ok {
			d.logger.Warn().
				Str("channel", string(msg.Channel)).
				Str("event_type", msg.EventType).
				Msg("no sender configured for channel")
			continue
		}

		if err := sender.Send(ctx, msg); err != nil {
			d.logger.Error().
				Err(err).
				Str("channel", string(msg.Channel)).
				Str("event_type", msg.EventType).
				Msg("failed to send notification")
			continue
		}

		if d.history != nil {
			d.history.Add(msg)
		}
	}

	return nil
}
