package queues

import (
	"github.com/PiotrFerenc/mash2/api/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	Connect() error
	AddStageToQueue(message types.Message) error
	AddStageAsFailed(message types.Message) error
	AddStageAsSuccess(message types.Message) error
	WaitingForFailedStage() (<-chan amqp.Delivery, error)
	WaitingForSucceedStage() (<-chan amqp.Delivery, error)
	WaitingForStage() (<-chan amqp.Delivery, error)
	CreateQueue(name string) error
}
