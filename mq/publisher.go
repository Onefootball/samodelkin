package mq

import (
	"fmt"
	"time"

	"github.com/motain/amqp"
)

type (
	// Publisher describes logic
	// for publishing messages to
	// the amqp server
	Publisher interface {
		Publish(b []byte, headers amqp.Table, routingKey string) error
	}

	// AMQPPublisher implements Publisher interface
	AMQPPublisher struct {
		amqpChannel *amqp.Channel
		exchange    string
	}
)

// NewAMQPPublisher inits and returns an AMQPublisher
func NewAMQPPublisher(amqpChannel *amqp.Channel, exchange string) *AMQPPublisher {
	return &AMQPPublisher{
		amqpChannel: amqpChannel,
		exchange:    exchange,
	}
}

// Publish sends a message to the amqp server
func (p *AMQPPublisher) Publish(b []byte, headers amqp.Table, routingKey string) error {
	pub := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Headers:      headers,
		Body:         b,
	}

	if err := p.amqpChannel.Publish(p.exchange, routingKey, true, false, pub); err != nil {
		return fmt.Errorf("AMQPPublisher.Publish: %s", err)
	}

	return nil
}
