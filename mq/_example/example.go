package main

import (
	"log"
	"os"

	"github.com/motain/amqp"
	"github.com/motain/samodelkin/mq"
)

const (
	ServiceLogPrefix = "[ARTHUR]: "
)

type (
	ExampleHandler struct {
		logger *log.Logger
	}
)

func NewExamplesHandler(logger *log.Logger) *ExampleHandler {
	return &ExampleHandler{db, logger}
}

func (h *ExampleHandler) HandleDelivery(d *amqp.Delivery) error {
	h.logger.Println("consuming queue")
	return nil
}

func main() {
	logger := log.New(os.Stderr, ServiceLogPrefix, log.LstdFlags)

	handlers := map[string]mq.AMQPHandler{
		"user.settings": NewExamplesHandler(logger),
	}

	cfgMQ := mq.RabbitMQConfig{
		Connection: mq.RabbitMQConnection{
			User:            "test",
			Pass:            "test",
			Host:            "localhost",
			Port:            5632,
			Attempts:        5,
			ErrLogEnable:    true,
			ReturnLogEnable: true,
		},
		Consumers: []mq.RabbitMQConsumer{
			mq.RabbitMQConsumer{
				ID:      1,
				Queue:   "test",
				Workers: 1,
			},
		},
	}

	amqpConnector := mq.NewAMQPConnector(cfgMQ, logger)
	amqpChannel, err := amqpConnector.Channel()
	if err != nil {
		logger.Panic(err)
	}

	err = amqpConnector.Consume(amqpChannel, handlers)
	if err != nil {
		logger.Panic(err)
	}

	select {}
}
