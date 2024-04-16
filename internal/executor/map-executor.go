package executor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions/common"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions/docker"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions/file"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions/git"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions/math"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions/others"
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

// CreateMapExecutor creates and starts an executor that consumes tasks from the given message queue and executes them.
// It takes a MessageQueue and a map of Actions as input parameters.
// The Actions map should contain the available actions indexed by their names.
// The function returns an Executor instance.
func CreateMapExecutor(queue queues.MessageQueue, actions map[string]actions.Action) Executor {
	a := actions

	go func() {
		Task, err := queue.WaitingForTask()
		if err != nil {
			log.Fatal(err)
		}
		var forever chan struct{}

		go func() {
			for d := range Task {
				message, err := unmarshal(d)

				action, ok := a[message.CurrentStep.Action]
				if ok {
					message, err = action.Execute(message)
					addToQueue(err, queue, message)
				} else {
					e := fmt.Sprintf("Action %s not found", message.CurrentStep.Action)
					addToQueue(errors.New(e), queue, message)
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

func addToQueue(err error, queue queues.MessageQueue, message types.Process) {
	if err != nil {
		err = queue.AddTaskAsFailed(err, message)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	err = queue.AddTaskAsSuccess(message)
	if err != nil {
		log.Fatal(err)
	}
}

func unmarshal(d amqp.Delivery) (types.Process, error) {
	var message types.Process
	err := json.Unmarshal(d.Body, &message)
	if err != nil {

		log.Fatal(err)
	}
	return message, err
}

func CreateActionMap(config *configuration.Config) map[string]actions.Action {
	return map[string]actions.Action{
		"console":     others.CreateConsoleAction(),
		"add-numbers": math.CreateAddNumbers(),
		"git-clone":   git.CreateGitClone(config),
		"git-commit":  git.CreateGitCommit(config),
		"git-branch":  git.CreateGitCreateBranch(config),
		"file-create": file.CreateContentToFile(config),
		"docker-run":  docker.CreateDockerRun(),
		"file-delete": file.CreateDeleteFileAction(config),
		"file-append": file.CreateAppendContentToFile(config),
		"for-each":    common.CreateForEachLoop(),
	}
}
