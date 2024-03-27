package services

import (
	"errors"
	apitypes "github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/queues"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/goccy/go-json"
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

func onSuccess(process types.Pipeline, queue queues.MessageQueue) error {
	index := process.CurrentStep.Sequence
	if index < len(process.Steps) {
		current := process.Steps[index]
		process.CurrentStep = current
		err := queue.AddStageToQueue(process)
		if err != nil {
			return err
		}
	} else {
		log.Printf("Done: ")
		return nil
	}

	return nil
}

func onFail(message types.Pipeline, _ queues.MessageQueue) error {
	log.Printf("Fail %s => %+v\\n", message.CurrentStep.Name, message.Error)
	return errors.New(message.Error)
}

func CreatePipelineService(queue queues.MessageQueue, processService ProcessService) PipelineService {
	go watchForMessages(queue, onSuccess, onFail, processService)

	return &pipelineService{
		queue:          queue,
		processService: processService,
	}
}

func watchForMessages(queue queues.MessageQueue, onSuccess func(types.Pipeline, queues.MessageQueue) error, onFail func(types.Pipeline, queues.MessageQueue) error, processService ProcessService) {
	successStages, _ := queue.WaitingForSucceedStage()
	failedStages, _ := queue.WaitingForFailedStage()
	var forever chan struct{}
	go func() {
		err := processStages(successStages, queue, onSuccess, processService)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := processStages(failedStages, queue, onFail, processService)
		if err != nil {
			panic(err)
		}
	}()
	<-forever
}

func processStages(stages <-chan amqp.Delivery, queue queues.MessageQueue, onFunc func(types.Pipeline, queues.MessageQueue) error, processService ProcessService) error {
	for d := range stages {
		message, err := unmarshalMessage(d.Body)
		if err != nil {
			processService.MarkAsFailed(message, err)
			return err
		}
		processService.MarkAsStarted(message)
		err = onFunc(*message, queue)
		if err != nil {
			processService.MarkAsFailed(message, err)
			return err
		}
		processService.MarkAsDone(message)
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

	err := p.queue.AddStageToQueue(*process)
	if err != nil {
		return err
	}

	return nil
}
