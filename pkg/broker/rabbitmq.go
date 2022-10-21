package broker

import (
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/proto"
	"context"
	"github.com/go-faster/errors"
	protobuf "github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
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
func (c RabbitMQBroker) ProcessDatabaseQuery(ctx context.Context, _ course.DepartmentID, msg *proto.CourseDatabaseBatchMessage) error {
	data, err := protobuf.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "cannot marshal")
	}
	return c.channel.PublishWithContext(ctx,
		"",
		c.queue.Name,
		false,
		false,
		amqp.Publishing{
			Body: data,
		})
}

// Consume will consume the messages which are received on
func (c RabbitMQBroker) Consume(consumer string) (<-chan *proto.CourseDatabaseBatchMessage, error) {
	// Create the consumer
	messages, err := c.channel.Consume(
		c.queue.Name,
		consumer,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	// Create the channel to send the parsed proto buffer messages in it
	messageChannel := make(chan *proto.CourseDatabaseBatchMessage)
	// Loop over messages and send them in another goroutine
	go receiveMessages(messages, messageChannel)
	return messageChannel, nil
}

// CancelConsumer will cancel a consumer by its name
func (c RabbitMQBroker) CancelConsumer(consumer string) error {
	return c.channel.Cancel(consumer, false)
}

// receiveMessages will receive messages from a channel and parses them as proto.CourseDatabaseBatchMessage
func receiveMessages(incoming <-chan amqp.Delivery, outgoing chan<- *proto.CourseDatabaseBatchMessage) {
	for message := range incoming {
		parsedMessage := new(proto.CourseDatabaseBatchMessage)
		err := protobuf.Unmarshal(message.Body, parsedMessage)
		if err != nil {
			log.WithError(err).Warn("cannot parse proto buffer message")
			continue
		}
		// Send to channel
		outgoing <- parsedMessage
	}
	// When incoming channel is closed, also close the outgoing channel
	close(outgoing)
}
