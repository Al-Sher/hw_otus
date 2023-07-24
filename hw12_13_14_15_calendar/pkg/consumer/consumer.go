package consumer

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type RMQConnection interface {
	Channel() (*amqp091.Channel, error)
}

type Consumer struct {
	name    string
	channel *amqp091.Channel
	conn    RMQConnection
}

type HandleFunc func(message []byte) error

func New(
	name string,
	conn RMQConnection,
	exchange string,
	exchangeType string,
	queueName string,
	routingKey string,
) (*Consumer, error) {
	var err error
	c := &Consumer{
		name:    name,
		conn:    conn,
		channel: nil,
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	return c, c.declare(exchange, exchangeType, queueName, routingKey)
}

type Message struct {
	Ctx  context.Context
	Data []byte
}

func (c *Consumer) Consume(ctx context.Context, queue string, handle HandleFunc) (<-chan Message, error) {
	messages := make(chan Message)

	go func() {
		<-ctx.Done()
		if err := c.channel.Close(); err != nil {
			return
		}
	}()

	deliveries, err := c.channel.Consume(queue, c.name, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(messages)

		for {
			select {
			case <-ctx.Done():
				return
			case delivery := <-deliveries:
				err = handle(delivery.Body)
				if err == nil {
					if err := delivery.Ack(false); err != nil {
						return
					}
				}
			}
		}
	}()

	return messages, nil
}

func (c *Consumer) declare(exchange string, exchangeType string, queueName string, routingKey string) error {
	if err := c.channel.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil); err != nil {
		return err
	}

	queue, err := c.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return c.channel.QueueBind(queue.Name, routingKey, exchange, false, nil)
}
