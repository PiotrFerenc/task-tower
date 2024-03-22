package queues

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/configuration"
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

	err = q.CreateQueue(configuration.QueueStageSucceed)
	if err != nil {
		panic(err)
	}

	err = q.CreateQueue(configuration.QueueStageFailed)
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
func (queue *queue) WaitingForFailedStage() (<-chan amqp.Delivery, error) {

	return queue.client.ch.Consume(
		queue.configuration.QueueStageFailed,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

}
func (queue *queue) WaitingForSucceedStage() (<-chan amqp.Delivery, error) {

	return queue.client.ch.Consume(
		queue.configuration.QueueStageSucceed,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

}
func (queue *queue) WaitingForStage() (<-chan amqp.Delivery, error) {

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
func (queue *queue) AddStageToQueue(message types.Message) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = queue.client.ch.PublishWithContext(ctx,
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

	return err
}
func (queue *queue) AddStageAsSuccess(message types.Message) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = queue.client.ch.PublishWithContext(ctx,
		"",                                    // exchange
		queue.configuration.QueueStageSucceed, // routing key
		false,                                 // mandatory
		false,                                 // immediate
		amqp.Publishing{
			ContentType:   ContentType,
			CorrelationId: uuid.NewString(),
			ReplyTo:       queue.configuration.QueueStageSucceed,
			Body:          bytes,
		})

	return err
}
func (queue *queue) AddStageAsFailed(message types.Message) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = queue.client.ch.PublishWithContext(ctx,
		"",
		queue.configuration.QueueStageFailed,
		false,
		false,
		amqp.Publishing{
			ContentType:   ContentType,
			CorrelationId: uuid.NewString(),
			ReplyTo:       queue.configuration.QueueStageFailed,
			Body:          bytes,
		})

	return err
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
