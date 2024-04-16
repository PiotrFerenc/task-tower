package services

import (
	"encoding/json"
	apitypes "github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/queues"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type PipelineService interface {
	Run(pipeline *apitypes.Pipeline) (uuid.UUID, error)
}

type pipelineService struct {
	queue          queues.MessageQueue
	processService ProcessService
}

func CreatePipelineService(queue queues.MessageQueue, processService ProcessService, onSuccess OnMessageReceivedFunc, onFail OnMessageReceivedFunc, onFinish OnMessageReceivedFunc) PipelineService {
	go watchForMessages(queue, onSuccess, onFail, onFinish, processService)

	return &pipelineService{
		queue:          queue,
		processService: processService,
	}
}
func (p *pipelineService) Run(pipeline *apitypes.Pipeline) (uuid.UUID, error) {
	process := types.NewProcessFromPipeline(pipeline)
	p.processService.MarkAsStarted(process)

	err := p.queue.AddTaskToQueue(*process)
	if err != nil {
		return uuid.Nil, err
	}

	return process.Id, nil
}
func watchForMessages(queue queues.MessageQueue, onSuccess OnMessageReceivedFunc, onFail OnMessageReceivedFunc, onFinish OnMessageReceivedFunc, processService ProcessService) {
	forever := make(chan struct{})
	go ProcessTasks(queue, onSuccess, processService, queue.WaitingForSucceedTask)
	go ProcessTasks(queue, onFail, processService, queue.WaitingForFailedTask)
	go ProcessTasks(queue, onFinish, processService, queue.WaitingForFinishedTask)
	<-forever
}

func ProcessTasks(queue queues.MessageQueue, taskFunc OnMessageReceivedFunc, processService ProcessService, getTasksFunc getTasksFunc) {
	tasks, _ := getTasksFunc()
	err := processTasks(tasks, queue, taskFunc, processService)
	if err != nil {
		panic(err)
	}
}

type getTasksFunc func() (<-chan amqp.Delivery, error)

func processTasks(tasks <-chan amqp.Delivery, queue queues.MessageQueue, taskFunc OnMessageReceivedFunc, processService ProcessService) error {
	for d := range tasks {
		message, err := unmarshalMessage(d.Body)
		if err != nil {
			return err
		}
		taskFunc(message, queue, processService)
	}
	return nil
}

func unmarshalMessage(body []byte) (*types.Process, error) {
	var message types.Process
	err := json.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

type OnMessageReceivedFunc func(process *types.Process, queue queues.MessageQueue, service ProcessService)

func CreateOnSuccessFunc() OnMessageReceivedFunc {
	return func(process *types.Process, queue queues.MessageQueue, service ProcessService) {
		if len(process.Steps) == 0 {
			err := queue.AddTaskAsFinished(*process)
			if err != nil {
				panic(err)
			}
		} else {
			currentStep := process.Steps[0]
			process.CurrentStep = currentStep
			process.Steps = process.Steps[1:]
			err := queue.AddTaskToQueue(*process)
			if err != nil {
				panic(err)
			}
			service.TaskFinished(process)
		}
	}

}

func CreateOnFailFunc() OnMessageReceivedFunc {
	return func(process *types.Process, queue queues.MessageQueue, service ProcessService) {
		log.Printf("Fail %s => %+v\\n", process.CurrentStep.Name, process.Error)
		service.MarkAsFailed(process, process.Error)
	}
}

func CreateOnFinishFunc() OnMessageReceivedFunc {
	return func(process *types.Process, queue queues.MessageQueue, service ProcessService) {
		log.Printf("Done %s ", process.Id)
		service.MarkAsDone(process)
	}
}
