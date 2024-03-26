package executor

import (
	"encoding/json"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions/file"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/queues"
	"github.com/PiotrFerenc/mash2/internal/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type MapExecutor struct {
}

type executor struct {
	queue queues.MessageQueue
}

func CreateMapExecutor(queue queues.MessageQueue, config *configuration.Config) Executor {
	a := map[string]actions.Action{
		"console":     actions.CreateConsoleAction(),
		"add-numbers": actions.CreateAddNumbers(),
		"git-clone":   actions.CreateGitClone(config),
		"file-create": file.CreateContentToFile(config),
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

				action, ok := a[message.CurrentStep.Action]
				if ok {
					message, err = action.Execute(message)
					addToQueue(err, queue, message)
				} else {
					log.Printf("No action: %s", message.CurrentStep.Action)
				}
			}
		}()

		log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
		<-forever
	}()

	return &executor{
		queue: queue,
	}
}

func addToQueue(err error, queue queues.MessageQueue, message types.Pipeline) {
	if err != nil {
		message.Error = err.Error()
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

func unmarshal(d amqp.Delivery) (types.Pipeline, error) {
	var message types.Pipeline
	err := json.Unmarshal(d.Body, &message)
	if err != nil {
		panic(err)
	}
	return message, err
}
