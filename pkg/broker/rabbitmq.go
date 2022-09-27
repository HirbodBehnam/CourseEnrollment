package broker

import (
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/proto"
	"context"
	"github.com/go-faster/errors"
	protobuf "github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQBroker instantiates a RabbitMQ broker for general use
type RabbitMQBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewRabbitMQBroker creates a RabbitMQ connection and declares a durable queue
func NewRabbitMQBroker(connectionUrl, queueName string) (RabbitMQBroker, error) {
	// Connect to rabbit mq server
	conn, err := amqp.Dial(connectionUrl)
	if err != nil {
		return RabbitMQBroker{}, err
	}
	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return RabbitMQBroker{}, err
	}
	// Create the queue
	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = conn.Close() // channel will be closed as well
		return RabbitMQBroker{}, err
	}
	return RabbitMQBroker{conn, ch, queue}, nil
}

func (c RabbitMQBroker) Close() error {
	// channel will be closed as well if conn is closed
	return c.conn.Close()
}

// ProcessDatabaseQuery will push a database query into queue.
// DepartmentID is currently unused.
func (c RabbitMQBroker) ProcessDatabaseQuery(_ course.DepartmentID, msg *proto.CourseDatabaseBatchMessage) error {
	data, err := protobuf.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "cannot marshal")
	}
	return c.channel.PublishWithContext(context.Background(),
		"",
		c.queue.Name,
		false,
		false,
		amqp.Publishing{
			Body: data,
		})
}
