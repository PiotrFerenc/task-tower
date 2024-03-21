package executor

import (
	"encoding/json"
	"github.com/PiotrFerenc/mash2/src/api/types"
	"github.com/PiotrFerenc/mash2/src/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/src/internal/queues"
	"log"
)

type MapExecutor struct {
}

type executor struct {
	queue queues.MessageQueue
}

func CreateMapExecutor(queue queues.MessageQueue) Executor {
	go func() {

		stage, err := queue.WaitingForStage()
		if err != nil {
		}
		var forever chan struct{}

		go func() {
			for d := range stage {
				var message types.Message
				err := json.Unmarshal(d.Body, &message)
				if err != nil {
					panic(err)
				}
				a := map[string]actions.Action{
					"hallo": actions.CreateHalloAction(),
				}
				action, _ := a[message.CurrentStage.Action]

				message, err = action.Execute(message)
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
		}()

		log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
		<-forever
	}()

	return &executor{
		queue: queue,
	}
}

func (executor *executor) Execute(actionName string, parameters types.Message) error {

	//TODO: remove
	return nil
}
