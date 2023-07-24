package producer

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RMQConnection interface {
	Channel() (*amqp091.Channel, error)
}

type Producer struct {
	exchange     string
	conn         RMQConnection
	exchangeType string
}

type HandleFunc func(message []byte) error

func New(exchange string, exchangeType string, conn RMQConnection) *Producer {
	return &Producer{
		exchange:     exchange,
		conn:         conn,
		exchangeType: exchangeType,
	}
}

type Message struct {
	Data []byte
}

func (p *Producer) Publish(ctx context.Context, routingKey string, body []byte) error {
	channel, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	if err := channel.ExchangeDeclare(p.exchange, p.exchangeType, true, false, false, false, nil); err != nil {
		return fmt.Errorf("exchange Declare: %w", err)
	}

	if err = channel.PublishWithContext(
		ctx,
		p.exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			Headers:         amqp091.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp091.Transient,
			Priority:        0,
		},
	); err != nil {
		return fmt.Errorf("exchange Publish: %w", err)
	}

	return nil
}
