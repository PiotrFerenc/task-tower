package queues

import (
	"context"
	"fmt"
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
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

	return q
}

func (queue *queue) CreateQueue(name string) error {
	return queue.client.CreateQueue(name, true, false)
}

func (queue *queue) Connect() error {

	conn, err := ConnectRabbitMQ(queue.configuration.QueueUser, queue.configuration.QueuePassword, queue.configuration.QueueHost, queue.configuration.QueuePort, queue.configuration.QueueVhost)
	if err != nil {
		panic(err)
	}

	client, err := NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}

	queue.client = &client

	return nil
}
func (queue *queue) Subscribe() {

	//TODO: implement
}
func (queue *queue) Publish(message types.Stage) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := queue.client.ch.PublishWithContext(ctx,
		"	",                              // exchange
		queue.configuration.QueueRunPipe, // routing key
		false,                            // mandatory
		false,                            // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: "1",
			ReplyTo:       queue.configuration.QueueRunPipe,
			Body:          []byte(message.Name),
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

func (rc RabbitClient) CreateQueue(name string, durable, autoDelete bool) error {
	_, err := rc.ch.QueueDeclare(name, durable, autoDelete, false, false, nil)
	return err
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}
