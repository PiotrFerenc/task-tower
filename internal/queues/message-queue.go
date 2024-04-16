package queues

import (
	"github.com/PiotrFerenc/mash2/internal/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	Connect() error
	AddTaskToQueue(message types.Process) error
	AddTaskAsFailed(error error, message types.Process) error
	AddTaskAsSuccess(message types.Process) error
	AddTaskAsFinished(message types.Process) error
	WaitingForFailedTask() (<-chan amqp.Delivery, error)
	WaitingForSucceedTask() (<-chan amqp.Delivery, error)
	WaitingForFinishedTask() (<-chan amqp.Delivery, error)
	WaitingForTask() (<-chan amqp.Delivery, error)
	CreateQueue(name string) error
}
