package services

import (
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

func onSuccess(process types.Process, queue queues.MessageQueue) {
	index := process.CurrentStep.Sequence
	if index < len(process.Steps) {
		current := process.Steps[index]
		process.CurrentStep = current
		err := queue.AddStageToQueue(process)
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("Done")
	}
}

func onFail(message types.Process, queue queues.MessageQueue) {
	log.Printf("Fail %s => %+v\\n", message.CurrentStep.Name, message.Error)
}

func CreatePipelineService(queue queues.MessageQueue, processService ProcessService) PipelineService {
	go watchForMessages(queue, onSuccess, onFail)
	return &pipelineService{
		queue:          queue,
		processService: processService,
	}
}

func watchForMessages(queue queues.MessageQueue, onSuccess func(types.Process, queues.MessageQueue), onFail func(types.Process, queues.MessageQueue)) {
	successStages, _ := queue.WaitingForSucceedStage()
	failedStages, _ := queue.WaitingForFailedStage()
	var forever chan struct{}
	go processStages(successStages, queue, onSuccess)
	go processStages(failedStages, queue, onFail)
	<-forever
}

func processStages(stages <-chan amqp.Delivery, queue queues.MessageQueue, onSuccess func(types.Process, queues.MessageQueue)) {
	for d := range stages {
		message, err := unmarshalMessage(d.Body)
		if err != nil {
			panic(err)
		}
		onSuccess(*message, queue)
	}
}

func unmarshalMessage(body []byte) (*types.Process, error) {
	var message types.Process
	err := json.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (p *pipelineService) Run(pipeline *apitypes.Pipeline) error {
	process := types.NewProcessFromPipeline(pipeline)

	err := p.queue.AddStageToQueue(*process)
	if err != nil {
		return err
	}

	return nil
}
