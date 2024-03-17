package executor

import (
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
				log.Printf(" [x] %s", d.Body)
				err := queue.AddStageAsSuccess(d.MessageId)
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
