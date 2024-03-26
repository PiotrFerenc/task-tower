package queues

import (
	"github.com/PiotrFerenc/mash2/internal/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	Connect() error
	AddStageToQueue(message types.Pipeline) error
	AddStageAsFailed(message types.Pipeline) error
	AddStageAsSuccess(message types.Pipeline) error
	WaitingForFailedStage() (<-chan amqp.Delivery, error)
	WaitingForSucceedStage() (<-chan amqp.Delivery, error)
	WaitingForStage() (<-chan amqp.Delivery, error)
	CreateQueue(name string) error
}
