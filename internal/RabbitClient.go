package internal

//
//import (
//	"context"
//	"github.com/google/uuid"
//	amqp "github.com/rabbitmq/amqp091-go"
//	"strconv"
//	"time"
//)
//
//type QueueProvider interface {
//	Connect() error
//	Controller(queue string, correlationId uuid.UUID, body []byte)
//}
//
//type RabbitClient struct {
//	Queue  amqp.Queue
//	Chanel *amqp.Channel
//}
//
//func (rabbit *RabbitClient) Publish(queue string, correlationId uuid.UUID, body []byte) {
//	msgs, err := rabbit.Chanel.Consume(
//		q.Name, // queues
//		"",     // consumer
//		true,   // auto-ack
//		false,  // exclusive
//		false,  // no-local
//		false,  // no-wait
//		nil,    // args
//	)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	corrId := correlationId.String()
//
//	err = rabbit.Chanel.PublishWithContext(ctx,
//		"",          // exchange
//		"rpc_queue", // routing key
//		false,       // mandatory
//		false,       // immediate
//		amqp.Publishing{
//			ContentType:   "text/plain",
//			CorrelationId: corrId,
//			ReplyTo:       rabbit.Queue.Name,
//			Body:          body,
//		})
//
//	for d := range msgs {
//		if corrId == d.CorrelationId {
//			res, err = strconv.Atoi(string(d.Body))
//			break
//		}
//	}
//
//}
//
//func (rabbit *RabbitClient) Connect() error {
//	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
//	if err != nil {
//		return err
//	}
//	defer func(conn *amqp.Connection) {
//		err := conn.Close()
//		if err != nil {
//
//		}
//	}(conn)
//
//	rabbit.Chanel, err = conn.Channel()
//	if err != nil {
//		return err
//	}
//	defer func(ch *amqp.Channel) {
//		err := ch.Close()
//		if err != nil {
//
//		}
//	}(rabbit.Chanel)
//
//	rabbit.Queue, err = rabbit.Chanel.QueueDeclare(
//		"",    // name
//		false, // durable
//		false, // delete when unused
//		true,  // exclusive
//		false, // noWait
//		nil,   // arguments
//	)
//
//	return nil
//}
