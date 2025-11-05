package queue

import (
	"context"
	"errors"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/youngermaster/distributed-bookstore/notification-service/internal/config"
	"github.com/youngermaster/distributed-bookstore/notification-service/internal/events"
	"github.com/youngermaster/distributed-bookstore/notification-service/internal/notification"
)

type Consumer struct {
	cfg         *config.Config
	logger      zerolog.Logger
	dispatcher  *notification.Dispatcher
	conn        *amqp091.Connection
	channel     *amqp091.Channel
	queue       amqp091.Queue
	consumerTag string
}

func NewConsumer(cfg *config.Config, dispatcher *notification.Dispatcher, logger zerolog.Logger) (*Consumer, error) {
	c := &Consumer{
		cfg:         cfg,
		logger:      logger,
		dispatcher:  dispatcher,
		consumerTag: cfg.ServiceName + "-consumer",
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Consumer) connect() error {
	conn, err := amqp091.Dial(c.cfg.RabbitMQURL)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}

	if err := channel.Qos(c.cfg.RabbitMQPrefetchCount, 0, false); err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}

	if err := channel.ExchangeDeclare(
		c.cfg.RabbitMQExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}

	queue, err := channel.QueueDeclare(
		c.cfg.RabbitMQQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}

	for _, key := range c.cfg.RabbitRoutingKeys {
		if err := channel.QueueBind(
			queue.Name,
			key,
			c.cfg.RabbitMQExchange,
			false,
			nil,
		); err != nil {
			_ = channel.Close()
			_ = conn.Close()
			return err
		}
	}

	c.conn = conn
	c.channel = channel
	c.queue = queue

	return nil
}

func (c *Consumer) Start(ctx context.Context) error {
	deliveries, err := c.channel.Consume(
		c.queue.Name,
		c.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.logger.Info().
		Str("queue", c.queue.Name).
		Strs("routing_keys", c.cfg.RabbitRoutingKeys).
		Msg("notification consumer started")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info().Msg("notification consumer stopping")
			return nil
		case delivery, ok := <-deliveries:
			if !ok {
				return errors.New("deliveries channel closed")
			}

			event, err := events.FromJSON(delivery.Body)
			if err != nil {
				c.logger.Error().
					Err(err).
					Msg("failed to decode notification event")
				_ = delivery.Nack(false, false)
				continue
			}

			if event.EventID == "" {
				event.EventID = delivery.MessageId
			}
			if event.OccurredAt.IsZero() {
				event.OccurredAt = time.Now()
			}

			if err := c.dispatcher.HandleEvent(ctx, event); err != nil {
				c.logger.Error().
					Err(err).
					Str("event_type", event.EventType).
					Msg("handler error processing event")
				_ = delivery.Nack(false, true)
				continue
			}

			if err := delivery.Ack(false); err != nil {
				c.logger.Error().Err(err).Msg("failed to ack message")
			}
		}
	}
}

func (c *Consumer) Close() {
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
