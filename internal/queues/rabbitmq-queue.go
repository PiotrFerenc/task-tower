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

func (q *queue) CreateQueue(name string) error {
	return q.client.CreateQueue(name)
}

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

func (q *queue) WaitingForFailedTask() (<-chan amqp.Delivery, error) {
	return q.waitingForTask(q.configuration.QueueTaskFailed)
}
func (q *queue) WaitingForSucceedTask() (<-chan amqp.Delivery, error) {
	return q.waitingForTask(q.configuration.QueueTaskSucceed)
}
func (q *queue) WaitingForTask() (<-chan amqp.Delivery, error) {
	return q.waitingForTask(q.configuration.QueueRunPipe)
}
func (q *queue) WaitingForFinishedTask() (<-chan amqp.Delivery, error) {
	return q.waitingForTask(q.configuration.QueueTaskFinished)
}

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

func (q *queue) AddTaskToQueue(message types.Process) error {
	return q.addTask(q.configuration.QueueRunPipe, message)
}
func (q *queue) AddTaskAsSuccess(message types.Process) error {
	return q.addTask(q.configuration.QueueTaskSucceed, message)
}
func (q *queue) AddTaskAsFinished(message types.Process) error {
	return q.addTask(q.configuration.QueueTaskFinished, message)
}
func (q *queue) AddTaskAsFailed(error error, message types.Process) error {
	return q.addTask(q.configuration.QueueTaskFailed, message)
}

type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func ConnectRabbitMQ(user, password, host, port, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost))
}

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

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}
