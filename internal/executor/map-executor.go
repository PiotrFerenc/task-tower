package executor

import (
	"encoding/json"
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/queues"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type MapExecutor struct {
}

type executor struct {
	queue queues.MessageQueue
}

func CreateMapExecutor(queue queues.MessageQueue) Executor {
	a := map[string]actions.Action{
		"hallo": actions.CreateHalloAction(),
	}

	go func() {
		stage, err := queue.WaitingForStage()
		if err != nil {
			log.Fatal(err)
		}
		var forever chan struct{}

		go func() {
			for d := range stage {
				message, err := unmarshal(d)

				action := a[message.CurrentStage.Action]

				message, err = action.Execute(message)
				addToQueue(err, queue, message)
			}
		}()

		log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
		<-forever
	}()

	return &executor{
		queue: queue,
	}
}

func addToQueue(err error, queue queues.MessageQueue, message types.Message) {
	if err != nil {
		err = queue.AddStageAsFailed(message)
		if err != nil {
			panic(err)
		}
	}
	err = queue.AddStageAsSuccess(message)
	if err != nil {
		panic(err)
	}
}

func unmarshal(d amqp.Delivery) (types.Message, error) {
	var message types.Message
	err := json.Unmarshal(d.Body, &message)
	if err != nil {
		panic(err)
	}
	return message, err
}

func (executor *executor) Execute(actionName string, parameters types.Message) error {

	//TODO: remove
	return nil
}
