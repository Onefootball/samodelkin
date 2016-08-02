package mq

import (
	"fmt"
	"log"

	"github.com/motain/amqp"
)

type (
	// AMQPHanlder describes logic for
	// processing an message from an amqp server
	AMQPHandler interface {
		HandleDelivery(*amqp.Delivery) error
	}

	// AMQPHandlerFunc implements AMQPHandler interface
	AMQPHandlerFunc func(*amqp.Delivery) error
)

// HandleDelivery calls f(d)
func (f AMQPHandlerFunc) HandleDelivery(d *amqp.Delivery) error {
	return f(d)
}

type (
	// Consumer interface desribes
	// logic for an amqp consumer
	Consumer interface {
		Consume() error
	}

	// AMQPConsumer implements Consumer interface
	AMQPConsumer struct {
		Handler     AMQPHandler
		AMQPChannel *amqp.Channel
		Queue       string
		Tag         string
		done        chan bool
		logger      *log.Logger
	}
)

// NewAMQPConsumer inits and returns a pointer
// to a new AMQPConsumer instance
func NewAMQPConsumer(h AMQPHandler, ch *amqp.Channel, logger *log.Logger, queue, tag string) *AMQPConsumer {
	return &AMQPConsumer{
		Handler:     h,
		AMQPChannel: ch,
		Queue:       queue,
		Tag:         tag,
		done:        make(chan bool),
		logger:      logger,
	}
}

// Consume processes amqp messages
func (c *AMQPConsumer) Consume() error {
	deliveries, err := c.AMQPChannel.Consume(c.Queue, c.Tag, false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func(ds <-chan amqp.Delivery) {
		for d := range ds {
			if err := tryCatchHandler(c.Handler, &d); err != nil {
				c.logger.Println(err)

				if err := d.Nack(false, false); err != nil {
					c.logger.Println(err)
				}

				continue
			}

			if err := d.Ack(false); err != nil {
				c.logger.Println(err)
			}
		}

		c.done <- true
	}(deliveries)

	return nil
}

// tryCatch calls h.HandleDelivery, catches a panic exception
// if any is thrown and returns it as an error
func tryCatchHandler(h AMQPHandler, d *amqp.Delivery) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %s", r)
		}
	}()
	return h.HandleDelivery(d)
}
