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

	err = q.CreateQueue(configuration.QueueRunPipe)
	if err != nil {
		panic(err)
	}

	err = q.CreateQueue(configuration.QueueTasksucceed)
	if err != nil {
		panic(err)
	}

	err = q.CreateQueue(configuration.QueueTaskFailed)
	if err != nil {
		panic(err)
	}

	err = q.CreateQueue(configuration.QueueTaskFinished)
	if err != nil {
		panic(err)
	}

	return q
}

func (queue *queue) CreateQueue(name string) error {
	return queue.client.CreateQueue(name)
}

func (queue *queue) Connect() error {
	var connection *amqp.Connection
	for {
		conn, err := ConnectRabbitMQ(queue.configuration.QueueUser, queue.configuration.QueuePassword, queue.configuration.QueueHost, queue.configuration.QueuePort, queue.configuration.QueueVhost)
		if err == nil {
			connection = conn
			log.Printf("Connected to rabbitmq")
			break
		} else {
			log.Printf("waiting for rabbitmq")
		}
		time.Sleep(1 * time.Second)
	}
	client, err := NewRabbitMQClient(connection)
	if err != nil {
		panic(err)
	}

	queue.client = &client

	return nil
}
func (queue *queue) WaitingForFailedTask() (<-chan amqp.Delivery, error) {

	return queue.client.ch.Consume(
		queue.configuration.QueueTaskFailed,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

}
func (queue *queue) WaitingForSucceedTask() (<-chan amqp.Delivery, error) {

	return queue.client.ch.Consume(
		queue.configuration.QueueTasksucceed,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

}
func (queue *queue) WaitingForTask() (<-chan amqp.Delivery, error) {

	return queue.client.ch.Consume(
		queue.configuration.QueueRunPipe,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

}
func (queue *queue) WaitingForFinishedTask() (<-chan amqp.Delivery, error) {

	return queue.client.ch.Consume(
		queue.configuration.QueueTaskFinished,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

}
func (queue *queue) AddTaskToQueue(message types.Process) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return queue.client.ch.PublishWithContext(ctx,
		"",
		queue.configuration.QueueRunPipe,
		false,
		false,
		amqp.Publishing{
			ContentType:   ContentType,
			CorrelationId: uuid.NewString(),
			ReplyTo:       queue.configuration.QueueRunPipe,
			Body:          bytes,
		})
}
func (queue *queue) AddTaskAsSuccess(message types.Process) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return queue.client.ch.PublishWithContext(ctx,
		"",                                   // exchange
		queue.configuration.QueueTasksucceed, // routing key
		false,                                // mandatory
		false,                                // immediate
		amqp.Publishing{
			ContentType:   ContentType,
			CorrelationId: uuid.NewString(),
			ReplyTo:       queue.configuration.QueueTasksucceed,
			Body:          bytes,
		})
}
func (queue *queue) AddTaskAsFinished(message types.Process) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return queue.client.ch.PublishWithContext(ctx,
		"",                                    // exchange
		queue.configuration.QueueTaskFinished, // routing key
		false,                                 // mandatory
		false,                                 // immediate
		amqp.Publishing{
			ContentType:   ContentType,
			CorrelationId: uuid.NewString(),
			ReplyTo:       queue.configuration.QueueTaskFinished,
			Body:          bytes,
		})

}
func (queue *queue) AddTaskAsFailed(error error, message types.Process) error {
	message.Error = error.Error()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return queue.client.ch.PublishWithContext(ctx,
		"",
		queue.configuration.QueueTaskFailed,
		false,
		false,
		amqp.Publishing{
			ContentType:   ContentType,
			CorrelationId: uuid.NewString(),
			ReplyTo:       queue.configuration.QueueTaskFailed,
			Body:          bytes,
		})
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
