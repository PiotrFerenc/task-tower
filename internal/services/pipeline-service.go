package services

import (
	"encoding/json"
	"errors"
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

func watchForMessages(queue queues.MessageQueue, onSuccess func(types.Pipeline, queues.MessageQueue) error, onFail func(types.Pipeline, queues.MessageQueue) error, onFinish func(types.Pipeline, queues.MessageQueue) error, processService ProcessService) {
	successStages, _ := queue.WaitingForSucceedStage()
	failedStages, _ := queue.WaitingForFailedStage()
	finishedStages, _ := queue.WaitingForFinishedStage()

	var forever chan struct{}

	go func() {
		message, err := processStages(finishedStages, queue, onFinish, processService)
		if err != nil {
			err := queue.AddStageAsFailed(err, *message)
			if err != nil {
				panic(err)
			}
		}
	}()

	go func() {
		message, err := processStages(successStages, queue, onSuccess, processService)
		if err != nil {
			err = queue.AddStageAsFailed(err, *message)
			if err != nil {
				panic(err)
			}
		}
	}()

	go func() {
		message, err := processStages(failedStages, queue, onFail, processService)
		if err != nil {
			err = queue.AddStageAsFailed(err, *message)
			if err != nil {
				panic(err)
			}
		}
	}()

	<-forever
}

func processStages(stages <-chan amqp.Delivery, queue queues.MessageQueue, onFunc func(types.Pipeline, queues.MessageQueue) error, processService ProcessService) (*types.Pipeline, error) {
	for d := range stages {
		message, err := unmarshalMessage(d.Body)
		if err != nil {
			processService.MarkAsFailed(message, err)
			return message, err
		}
		processService.MarkAsStarted(message)
		err = onFunc(*message, queue)
		if err != nil {
			processService.MarkAsFailed(message, err)
			return message, err
		}
		processService.MarkAsDone(message)
	}
	return nil, nil
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

	err := p.queue.AddStageToQueue(*process)
	if err != nil {
		return err
	}

	return nil
}

type OnSuccessFunc func(process types.Pipeline, queue queues.MessageQueue) error
type OnFailFunc func(process types.Pipeline, queue queues.MessageQueue) error
type OnFinishFunc func(process types.Pipeline, queue queues.MessageQueue) error

func CreateOnSuccessFunc() OnSuccessFunc {
	return func(process types.Pipeline, queue queues.MessageQueue) error {
		index := process.CurrentStep.Sequence
		if index < len(process.Steps) {
			current := process.Steps[index]
			process.CurrentStep = current
			err := queue.AddStageToQueue(process)
			if err != nil {
				return err
			}
		} else {
			err := queue.AddStageAsFinished(process)
			if err != nil {
				return err
			}
		}
		return nil
	}

}

func CreateOnFailFunc() OnFailFunc {
	return func(process types.Pipeline, queue queues.MessageQueue) error {
		log.Printf("Fail %s => %+v\\n", process.CurrentStep.Name, process.Error)
		return errors.New(process.Error)
	}
}

func CreateOnFinishFunc() OnFinishFunc {
	return func(process types.Pipeline, queue queues.MessageQueue) error {
		log.Printf("Done %s ", process.Id)
		return nil
	}
}
