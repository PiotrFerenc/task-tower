package queues

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	amqp "github.com/rabbitmq/amqp091-go"
)

type queue struct {
	configuration *configuration.QueueConfig
	client        RabbitClient
}

func CreateRabbitMqMessageQueue(configuration *configuration.QueueConfig) MessageQueue {
	return &queue{
		configuration: configuration,
	}
}

func (queue *queue) Connect() error {

	conn, err := ConnectRabbitMQ(queue.configuration.QueueUser, queue.configuration.QueuePassword, queue.configuration.QueueHost, queue.configuration.QueuePort, queue.configuration.QueueVhost)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client, err := NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	queue.client = client

	return nil
}
func (queue *queue) Receive() {
	//TODO: implement
}
func (queue *queue) Publish(message types.Stage) error {
	//TODO: implement
	return nil
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

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}
