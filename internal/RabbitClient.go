package internal

import (
	"context"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type QueueProvider interface {
	Connect() error
	Publish(queue string, correlationId uuid.UUID, body []byte)
}

type RabbitClient struct {
	Queue  amqp.Queue
	Ctx    context.Context
	Chanel *amqp.Channel
}

func (rabbit *RabbitClient) Publish(queue string, correlationId uuid.UUID, body []byte) {
	err := rabbit.Chanel.PublishWithContext(rabbit.Ctx,
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: correlationId.String(),
			ReplyTo:       rabbit.Queue.Name,
			Body:          body,
		})

	if err != nil {
		panic(err)
	}

}

func (rabbit *RabbitClient) Connect() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	rabbit.Chanel, err = conn.Channel()
	if err != nil {
		return err
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {

		}
	}(rabbit.Chanel)

	rabbit.Queue, err = rabbit.Chanel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rabbit.Ctx = ctx

	return nil
}
