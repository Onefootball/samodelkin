package mq

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/motain/amqp"
)

type (
	// RabbitMQConfig holds rabbitmq server configs
	RabbitMQConfig struct {
		Connection RabbitMQConnection `json:"connection"`
		Consumers  []RabbitMQConsumer `json:"consumers"`
	}

	// RabbitMQConnection desribes a single
	// rabbitmq server configuration
	RabbitMQConnection struct {
		User            string `json:"user"`
		Pass            string `json:"pass"`
		Host            string `json:"host"`
		Port            int    `json:"port"`
		Attempts        int    `json:"attempts"`
		ErrLogEnable    bool   `json:"error_log_enable"`
		ReturnLogEnable bool   `json:"return_log_enable"`
	}

	// RabbitMQConsumer desribes a rabbitmq consumer config
	RabbitMQConsumer struct {
		ID      string `json:"id"`
		Queue   string `json:"queue"`
		Workers int    `json:"workers"`
	}
)

// String returns a rabbitmq server connection string
func (c *RabbitMQConnection) String() string {
	return "amqp://" + c.User + ":" + c.Pass + "@" + c.Host + ":" + strconv.Itoa(c.Port)
}

type (
	// Connector describes amqp dial logic
	Connector interface {
		Channel() (*amqp.Channel, error)
		Consume(*amqp.Channel, map[string]AMQPHandler) error
	}

	// AMQPConnector implements Connector interface
	AMQPConnector struct {
		rabbitConfig *RabbitMQConfig
		logger       *log.Logger
	}
)

func NewAMQPConnector(cfg *RabbitMQConfig, logger *log.Logger) *AMQPConnector {
	return &AMQPConnector{cfg, logger}
}

func (f AMQPConnector) Channel() (*amqp.Channel, error) {
	_, amqpChannel, err := f.dial()
	if err != nil {
		return nil, err
	}

	if f.rabbitConfig.Connection.ErrLogEnable {
		go func() {
			f.logger.Println(fmt.Sprintf("closing: %s", <-amqpChannel.NotifyClose(make(chan *amqp.Error))))
		}()
	}

	if f.rabbitConfig.Connection.ReturnLogEnable {
		go func() {
			amqpReturns := amqpChannel.NotifyReturn(make(chan amqp.Return))
			for amqpR := range amqpReturns {
				f.logger.Printf("non-deliverable message: %s; routing key: %s", amqpR.Body, amqpR.RoutingKey)
			}
		}()
	}

	return amqpChannel, nil
}

func (f AMQPConnector) Consume(amqpChannel *amqp.Channel, handlers map[string]AMQPHandler) error {
	for _, consumerCfg := range f.rabbitConfig.Consumers {
		mqHandler, ok := handlers[consumerCfg.ID]
		if !ok {
			return fmt.Errorf("handler for id: %s not found!", consumerCfg.ID)
		}

		for i := 0; i < consumerCfg.Workers; i++ {
			consumer := NewAMQPConsumer(
				mqHandler,
				amqpChannel,
				f.logger,
				consumerCfg.Queue,
				fmt.Sprintf("%s_%d", consumerCfg.ID, i),
			)

			if err := consumer.Consume(); err != nil {
				return err
			}

			f.logger.Println(fmt.Sprintf("consuming from queue: %s", consumerCfg.Queue))
		}
	}

	return nil
}

// dial makes n attempts to dial an amqp
// server by a provided URL
func (f AMQPConnector) dial() (c *amqp.Connection, ch *amqp.Channel, err error) {
	if f.rabbitConfig.Connection.Attempts <= 0 {
		return nil, nil, errors.New("'attempts' must greater than 0")
	}

	for i := 0; i < f.rabbitConfig.Connection.Attempts; i++ {
		c, ch, err = f.dialOnce()
		if err == nil {
			return
		}

		time.Sleep(time.Second)
	}

	return
}

// dialAMQPOnce dials an amqp server under
// the provided URL
func (f AMQPConnector) dialOnce() (*amqp.Connection, *amqp.Channel, error) {
	c, err := amqp.Dial(f.rabbitConfig.Connection.String())
	if err != nil {
		return nil, nil, err
	}

	ch, err := c.Channel()
	if err != nil {
		c.Close()
		return nil, nil, err
	}

	return c, ch, nil
}
