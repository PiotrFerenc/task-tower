package queues

import (
	"github.com/PiotrFerenc/mash2/internal/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	Connect() error
	AddTaskToQueue(message types.Pipeline) error
	AddTaskAsFailed(error error, message types.Pipeline) error
	AddTaskAsSuccess(message types.Pipeline) error
	AddTaskAsFinished(message types.Pipeline) error
	WaitingForFailedTask() (<-chan amqp.Delivery, error)
	WaitingForSucceedTask() (<-chan amqp.Delivery, error)
	WaitingForFinishedTask() (<-chan amqp.Delivery, error)
	WaitingForTask() (<-chan amqp.Delivery, error)
	CreateQueue(name string) error
}
