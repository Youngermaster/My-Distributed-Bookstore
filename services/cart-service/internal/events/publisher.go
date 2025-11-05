package events

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	exchange string
	conn     *amqp091.Connection
	channel  *amqp091.Channel
}

func NewPublisher(url, exchange string) (*Publisher, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	if err := ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}

	return &Publisher{
		exchange: exchange,
		conn:     conn,
		channel:  ch,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, routingKey string, event NotificationEvent) error {
	body, err := event.Marshal()
	if err != nil {
		return err
	}

	err = p.channel.PublishWithContext(ctx, p.exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
		MessageId:   event.EventID,
		Timestamp:   event.OccurredAt,
		Type:        event.EventType,
		AppId:       event.Source,
	})
	if err != nil {
		return err
	}

	log.Printf("Published cart event %s to %s", event.EventType, routingKey)
	return nil
}

func (p *Publisher) Close() {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}
