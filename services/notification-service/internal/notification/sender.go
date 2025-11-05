package notification

import (
	"context"

	"github.com/rs/zerolog"
)

type Sender interface {
	Channel() Channel
	Send(ctx context.Context, msg Message) error
}

type logSender struct {
	channel Channel
	logger  zerolog.Logger
}

func NewLogSender(channel Channel, logger zerolog.Logger) Sender {
	return &logSender{
		channel: channel,
		logger:  logger,
	}
}

func (s *logSender) Channel() Channel {
	return s.channel
}

func (s *logSender) Send(ctx context.Context, msg Message) error {
	s.logger.Info().
		Str("channel", string(msg.Channel)).
		Str("event", msg.EventType).
		Str("recipient", msg.Recipient).
		Fields(msg.Metadata).
		Msg(msg.Body)
	return nil
}
