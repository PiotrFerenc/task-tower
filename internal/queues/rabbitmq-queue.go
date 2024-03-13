package queues

import (
	"github.com/PiotrFerenc/mash2/api/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type queue struct {
}

func CreateRabbitMqMessageQueue() MessageQueue {
	return &queue{}
}

func (queue *queue) Connect() error {
	conn, err := amqp.Dial("")
	if err != nil {
		return err
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	//rabbit.Chanel, err = conn.Channel()
	//if err != nil {
	//	return err
	//}
	//defer func(ch *amqp.Channel) {
	//	err := ch.Close()
	//	if err != nil {
	//
	//	}
	//}(rabbit.Chanel)

	return nil
}
func (queue *queue) Receive() {
	//TODO: implement
}
func (queue *queue) Publish(message types.Stage) error {
	//TODO: implement
	return nil
}
