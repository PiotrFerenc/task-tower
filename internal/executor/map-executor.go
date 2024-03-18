package executor

import (
	"encoding/json"
	"fmt"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/queues"
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
				log.Printf(" execute [x] %s", d.MessageId)
				var message queues.Message
				err := json.Unmarshal(d.Body, &message)
				if err != nil {
					panic(err)
				}
				a := map[string]actions.Action{
					"hallo": actions.Hallo{},
					"sleep": actions.Sleep{},
				}
				action, _ := a[message.CurrentStage.Name]

				action.Execute(actions.ActionContext{Parameters: message.CurrentStage.Parameters})

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

func (executor *executor) Execute(actionName string, parameters actions.ActionContext) error {

	a := map[string]actions.Action{
		"hallo": actions.Hallo{},
		"sleep": actions.Sleep{},
	}
	action, exist := a[actionName]
	if !exist {
		return fmt.Errorf("action %v not found", actionName)
	}
	action.Execute(parameters)

	return nil
}
