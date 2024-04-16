package services

import (
	"encoding/json"
	apitypes "github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/queues"
	"github.com/PiotrFerenc/mash2/internal/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type PipelineService interface {
	Run(pipeline *apitypes.Pipeline) error
}

type pipelineService struct {
	queue          queues.MessageQueue
	processService ProcessService
}

func CreatePipelineService(queue queues.MessageQueue, processService ProcessService, onSuccess OnSuccessFunc, onFail OnFailFunc, onFinish OnFinishFunc) PipelineService {
	go watchForMessages(queue, onSuccess, onFail, onFinish, processService)

	return &pipelineService{
		queue:          queue,
		processService: processService,
	}
}

func watchForMessages(queue queues.MessageQueue, onSuccess func(*types.Pipeline, queues.MessageQueue, ProcessService) error, onFail func(*types.Pipeline, queues.MessageQueue, ProcessService), onFinish func(*types.Pipeline, queues.MessageQueue, ProcessService), processService ProcessService) {
	successTasks, _ := queue.WaitingForSucceedTask()
	failedTasks, _ := queue.WaitingForFailedTask()
	finishedTasks, _ := queue.WaitingForFinishedTask()

	var forever chan struct{}

	go func() {
		err := processFinishTasks(finishedTasks, queue, onFinish, processService)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := processSuccessTasks(successTasks, queue, onSuccess, processService)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := processFailTasks(failedTasks, queue, onFail, processService)
		if err != nil {
			panic(err)
		}
	}()

	<-forever
}

func processSuccessTasks(Tasks <-chan amqp.Delivery, queue queues.MessageQueue, onFunc func(*types.Pipeline, queues.MessageQueue, ProcessService) error, processService ProcessService) error {
	for d := range Tasks {
		message, err := unmarshalMessage(d.Body)
		if err != nil {
			return err
		}
		err = onFunc(message, queue, processService)
		if err != nil {
			return err
		}
	}
	return nil
}
func processFinishTasks(Tasks <-chan amqp.Delivery, queue queues.MessageQueue, onFunc func(*types.Pipeline, queues.MessageQueue, ProcessService), processService ProcessService) error {
	for d := range Tasks {
		message, err := unmarshalMessage(d.Body)
		if err != nil {
			return err
		}
		onFunc(message, queue, processService)
	}
	return nil
}
func processFailTasks(Tasks <-chan amqp.Delivery, queue queues.MessageQueue, onFunc func(*types.Pipeline, queues.MessageQueue, ProcessService), processService ProcessService) error {
	for d := range Tasks {
		message, err := unmarshalMessage(d.Body)
		if err != nil {
			return err
		}
		onFunc(message, queue, processService)
	}
	return nil
}

func unmarshalMessage(body []byte) (*types.Pipeline, error) {
	var message types.Pipeline
	err := json.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (p *pipelineService) Run(pipeline *apitypes.Pipeline) error {
	process := types.NewProcessFromPipeline(pipeline)
	p.processService.MarkAsStarted(process)

	err := p.queue.AddTaskToQueue(*process)
	if err != nil {
		return err
	}

	return nil
}

type OnSuccessFunc func(process *types.Pipeline, queue queues.MessageQueue, service ProcessService) error
type OnFailFunc func(process *types.Pipeline, queue queues.MessageQueue, service ProcessService)
type OnFinishFunc func(process *types.Pipeline, queue queues.MessageQueue, service ProcessService)

func CreateOnSuccessFunc() OnSuccessFunc {
	return func(process *types.Pipeline, queue queues.MessageQueue, service ProcessService) error {
		if len(process.Steps) == 0 {
			err := queue.AddTaskAsFinished(*process)
			if err != nil {
				return err
			}
		} else {
			currentStep := process.Steps[0]
			process.CurrentStep = currentStep
			process.Steps = process.Steps[1:]
			err := queue.AddTaskToQueue(*process)
			if err != nil {
				return err
			}
			service.TaskFinished(process)
		}

		return nil
	}

}

func CreateOnFailFunc() OnFailFunc {
	return func(process *types.Pipeline, queue queues.MessageQueue, service ProcessService) {
		log.Printf("Fail %s => %+v\\n", process.CurrentStep.Name, process.Error)
		service.MarkAsFailed(process, process.Error)
	}
}

func CreateOnFinishFunc() OnFinishFunc {
	return func(process *types.Pipeline, queue queues.MessageQueue, service ProcessService) {
		log.Printf("Done %s ", process.Id)
		service.MarkAsDone(process)
	}
}
