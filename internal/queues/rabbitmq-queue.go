package queues

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

const (
	ContentType = "text/plain"
)

type queue struct {
	configuration *configuration.QueueConfig
	client        *RabbitClient
}

// CreateRabbitMqMessageQueue creates a new RabbitMQ message queue based on the given configuration.
// It establishes a connection, creates the necessary queues, and returns the created queue.
func CreateRabbitMqMessageQueue(configuration *configuration.QueueConfig) MessageQueue {
	q := &queue{
		configuration: configuration,
	}
	err := q.Connect()
	if err != nil {
		panic(err)
	}
	err = createQueues(q, configuration)
	if err != nil {
		panic(err)
	}
	return q
}

// createQueues creates the necessary queues based on the provided configuration.
// It iterates over the list of queue names in the configuration and calls the CreateQueue method
// to create each queue using the RabbitMQ client.
// Returns an error if any of the queue creations fail.
func createQueues(q *queue, configuration *configuration.QueueConfig) error {
	queues := []string{
		configuration.QueueRunPipe,
		configuration.QueueTaskSucceed,
		configuration.QueueTaskFailed,
		configuration.QueueTaskFinished,
	}
	for _, queueName := range queues {
		if err := q.CreateQueue(queueName); err != nil {
			return err
		}
	}
	return nil
}

// CreateQueue creates a new queue with the specified name using the RabbitMQ client.
// It uses the underlying RabbitClient's CreateQueue method to perform the creation.
// Returns an error if the queue creation fails.
func (q *queue) CreateQueue(name string) error {
	return q.client.CreateQueue(name)
}

// Connect establishes a connection to the RabbitMQ server and initializes a RabbitClient.
// It uses the underlying establishConnection function to establish the connection,
// and then creates a new RabbitMQ client using the connection.
// The resulting client is stored in the queue's client field for future use.
// Returns an error if the connection or client creation fails.
func (q *queue) Connect() error {
	connection, err := q.establishConnection()
	if err != nil {
		panic(err)
	}

	client, err := NewRabbitMQClient(connection)
	if err != nil {
		panic(err)
	}

	q.client = &client

	return nil
}

// establishConnection establishes a connection to the RabbitMQ server using the provided
// credentials and creates a RabbitClient with the connection.
// It retries the connection every 1 second until a successful connection is established.
// Returns the established connection and any error that occurred during the connection process.
func (q *queue) establishConnection() (conn *amqp.Connection, err error) {
	for {
		conn, err = ConnectRabbitMQ(q.configuration.QueueUser, q.configuration.QueuePassword, q.configuration.QueueHost, q.configuration.QueuePort, q.configuration.QueueVhost)
		if err == nil {
			log.Printf("Connected to rabbitmq")
			break
		} else {
			log.Printf("waiting for rabbitmq")
			time.Sleep(1 * time.Second)
		}
	}
	return conn, nil
}

// WaitingForFailedTask returns a channel of amqp.Delivery that receives messages from the `QueueTaskFailed` queue.
// It uses the `waitingForTask` method with the `QueueTaskFailed` queue name to retrieve the channel.
// Returns the channel and any error that occurred during the process.
func (q *queue) WaitingForFailedTask() (<-chan amqp.Delivery, error) {
	return q.waitingForTask(q.configuration.QueueTaskFailed)
}

// WaitingForSucceedTask returns a channel of amqp.Delivery that receives messages from
// the `QueueTaskSucceed` queue. It uses the `waitingForTask` method with the `QueueTaskSucceed`
// queue name to retrieve the channel. Returns the channel and any error that occurred during the process.
func (q *queue) WaitingForSucceedTask() (<-chan amqp.Delivery, error) {
	return q.waitingForTask(q.configuration.QueueTaskSucceed)
}

// WaitingForTask returns a channel of amqp.Delivery that receives messages from the
// `QueueRunPipe` queue. It uses the `waitingForTask` method with the `QueueRunPipe`
// queue name to retrieve the channel. Returns the channel and any error
// that occurred during the process.
func (q *queue) WaitingForTask() (<-chan amqp.Delivery, error) {
	return q.waitingForTask(q.configuration.QueueRunPipe)
}

// WaitingForFinishedTask returns a channel of amqp.Delivery that receives messages from the `QueueTaskFinished` queue.
// It uses the `waitingForTask` method with the `QueueTaskFinished` queue name to retrieve the channel.
// Returns the channel and any error that occurred during the process.
func (q *queue) WaitingForFinishedTask() (<-chan amqp.Delivery, error) {
	return q.waitingForTask(q.configuration.QueueTaskFinished)
}

// `waitingForTask` returns a channel of amqp.Delivery that receives messages from the specified queue.
// It uses the RabbitMQ client's `Consume` method to retrieve the channel.
// The function sets the arguments for `Consume` and passes them to the client.
// Returns the channel and any error that occurred during the process.
func (q *queue) waitingForTask(queueName string) (<-chan amqp.Delivery, error) {
	return q.client.ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

// addTask publishes a message to the specified queue using the RabbitMQ client.
// It marshals the provided message to JSON format and sets the message properties.
// Returns an error if the task addition fails.
func (q *queue) addTask(queueName string, message types.Process) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return q.client.ch.PublishWithContext(ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:   ContentType,
			CorrelationId: uuid.NewString(),
			ReplyTo:       queueName,
			Body:          bytes,
		})
}

// AddTaskToQueue adds a task to the specified queue using the RabbitMQ client.
// It calls the addTask method with the configuration's QueueRunPipe queue name and the provided message.
// Returns an error if the task addition fails.
func (q *queue) AddTaskToQueue(message types.Process) error {
	return q.addTask(q.configuration.QueueRunPipe, message)
}

// AddTaskAsSuccess adds a task to the specified queue as a successful task using the RabbitMQ client.
// It calls the addTask method with the configuration's QueueTaskSucceed queue name and the provided message.
// Returns an error if the task addition fails.
func (q *queue) AddTaskAsSuccess(message types.Process) error {
	return q.addTask(q.configuration.QueueTaskSucceed, message)
}

// AddTaskAsFinished adds a task to the specified queue as a finished task using the RabbitMQ client.
// It calls the addTask method with the configuration's QueueTaskFinished queue name and the provided message.
// Returns an error if the task addition fails.
func (q *queue) AddTaskAsFinished(message types.Process) error {
	return q.addTask(q.configuration.QueueTaskFinished, message)
}

// AddTaskAsFailed addTask adds a task to the specified queue using the RabbitMQ client.
// It marshals the message to JSON format and publishes it to the queue.
// Returns an error if the task addition fails.
func (q *queue) AddTaskAsFailed(error error, message types.Process) error {
	message.Error = error.Error()
	return q.addTask(q.configuration.QueueTaskFailed, message)
}

// RabbitClient represents a client for interacting with RabbitMQ.
//
// conn: Represents RabbitMQ connection.
// ch: Represents RabbitMQ channel.
type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// ConnectRabbitMQ establishes a connection to the RabbitMQ server using the provided
// credentials and returns the connection and any error that occurred during the connection process.
func ConnectRabbitMQ(user, password, host, port, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost))
}

// NewRabbitMQClient creates a new RabbitMQ client based on the provided connection.
// It establishes a channel using the connection and returns the created RabbitClient.
// Returns the created RabbitClient and any error that occurred during the process.
func NewRabbitMQClient(connection *amqp.Connection) (RabbitClient, error) {
	ch, err := connection.Channel()
	if err != nil {
		return RabbitClient{}, err
	}
	return RabbitClient{
		conn: connection,
		ch:   ch,
	}, nil
}

// CreateQueue creates a new queue with the specified name using the RabbitMQ client.
// It uses the underlying RabbitClient's QueueDeclare method to perform the creation.
// Returns an error if the queue creation fails.
func (rc RabbitClient) CreateQueue(name string) error {
	_, err := rc.ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	return err
}

// Close closes the RabbitMQ channel associated with the RabbitClient instance.
// It uses the underlying RabbitMQ channel's Close() method to close the channel.
// Returns an error if there is an issue closing the channel.
func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}
