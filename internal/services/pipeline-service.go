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
	successStages, _ := queue.WaitingForSucceedStage()
	failedStages, _ := queue.WaitingForFailedStage()
	finishedStages, _ := queue.WaitingForFinishedStage()

	var forever chan struct{}

	go func() {
		err := processFinishStages(finishedStages, queue, onFinish, processService)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := processSuccessStages(successStages, queue, onSuccess, processService)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := processFailStages(failedStages, queue, onFail, processService)
		if err != nil {
			panic(err)
		}
	}()

	<-forever
}

func processSuccessStages(stages <-chan amqp.Delivery, queue queues.MessageQueue, onFunc func(*types.Pipeline, queues.MessageQueue, ProcessService) error, processService ProcessService) error {
	for d := range stages {
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
func processFinishStages(stages <-chan amqp.Delivery, queue queues.MessageQueue, onFunc func(*types.Pipeline, queues.MessageQueue, ProcessService), processService ProcessService) error {
	for d := range stages {
		message, err := unmarshalMessage(d.Body)
		if err != nil {
			return err
		}
		onFunc(message, queue, processService)
	}
	return nil
}
func processFailStages(stages <-chan amqp.Delivery, queue queues.MessageQueue, onFunc func(*types.Pipeline, queues.MessageQueue, ProcessService), processService ProcessService) error {
	for d := range stages {
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

	err := p.queue.AddStageToQueue(*process)
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
		service.StageFinished(process)
		index := process.CurrentStep.Sequence
		if index < len(process.Steps) {
			current := process.Steps[index]
			process.CurrentStep = current
			err := queue.AddStageToQueue(*process)
			if err != nil {
				return err
			}
			service.StageFinished(process)
		} else {
			err := queue.AddStageAsFinished(*process)
			if err != nil {
				return err
			}
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
