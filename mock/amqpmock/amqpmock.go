// Package amqpmock provides basic data structures
// for mocking amqp related message ack/nack/reject logic
package amqpmock

import (
	"github.com/streadway/amqp"
)

type MockAcknowledger struct {
	amqp.Acknowledger
	AckFunc    func(tag uint64, multiple bool) error
	NackFunc   func(tag uint64, multiple bool, requeue bool) error
	RejectFunc func(tag uint64, requeue bool) error
}

// Ack calls m.AckFunc(tag, multiple).
//
// Ack signals an amqp server server that
// a message has been successfully processed
func (m *MockAcknowledger) Ack(tag uint64, multiple bool) error {
	return m.AckFunc(tag, multiple)
}

// Nack calls m.NackFunc(tag, multiple, requeue)
//
// Nack signals a negative acknowledgement(messages has been successfully processed)
// to the mq server.
//
// In contrast to Reject, Nack supports bulk (multiple) acknowledgement
func (m *MockAcknowledger) Nack(tag uint64, multiple bool, requeue bool) error {
	return m.NackFunc(tag, multiple, requeue)
}

// Reject calls m.RejectFunc(tag, requeue)
//
// Nack signals a negative acknowledgement(messages has been successfully processed)
// to the mq server.
//
// In contrast to Nack, Reject does not support bulk (multiple) acknowledgement
func (m *MockAcknowledger) Reject(tag uint64, requeue bool) error {
	return m.RejectFunc(tag, requeue)
}
