package queues

import (
	"github.com/PiotrFerenc/mash2/api/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	Connect() error
	AddStageToQueue(message Message) error
	AddStageAsFailed(message Message) error
	AddStageAsSuccess(message Message) error
	WaitingForFailedStage() (<-chan amqp.Delivery, error)
	WaitingForSucceedStage() (<-chan amqp.Delivery, error)
	WaitingForStage() (<-chan amqp.Delivery, error)
	CreateQueue(name string) error
}

type Message struct {
	CurrentStage types.Stage
	Pipeline     types.Pipeline
}
