package queues

import (
	"github.com/PiotrFerenc/mash2/internal/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	Connect() error
	AddStageToQueue(message types.Pipeline) error
	AddStageAsFailed(error error, message types.Pipeline) error
	AddStageAsSuccess(message types.Pipeline) error
	AddStageAsFinished(message types.Pipeline) error
	WaitingForFailedStage() (<-chan amqp.Delivery, error)
	WaitingForSucceedStage() (<-chan amqp.Delivery, error)
	WaitingForFinishedStage() (<-chan amqp.Delivery, error)
	WaitingForStage() (<-chan amqp.Delivery, error)
	CreateQueue(name string) error
}
