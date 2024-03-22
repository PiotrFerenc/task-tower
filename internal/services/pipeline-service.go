package services

import (
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/queues"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type PipelineService interface {
	Run(pipeline *types.Pipeline) error
}

type pipelineService struct {
	queue          queues.MessageQueue
	processService ProcessService
}

func onSuccess(message types.Message, queue queues.MessageQueue) {
	index := message.CurrentStage.Order
	if index < len(message.Pipeline.Stages) {
		current := message.Pipeline.Stages[index]
		err := queue.AddStageToQueue(types.Message{
			CurrentStage: current,
			Pipeline:     message.Pipeline,
		})
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("Done")
	}
}

func onFail(message types.Message, queue queues.MessageQueue) {
	log.Printf(" fail [x] %s", message.CurrentStage.Name)
}

func CreatePipelineService(queue queues.MessageQueue, processService ProcessService) PipelineService {
	go watchForMessages(queue, onSuccess, onFail)
	return &pipelineService{
		queue:          queue,
		processService: processService,
	}
}

func watchForMessages(queue queues.MessageQueue, onSuccess func(types.Message, queues.MessageQueue), onFail func(types.Message, queues.MessageQueue)) {
	successStages, _ := queue.WaitingForSucceedStage()
	failedStages, _ := queue.WaitingForFailedStage()
	var forever chan struct{}
	go processStages(successStages, queue, onSuccess)
	go processStages(failedStages, queue, onFail)
	<-forever
}

func processStages(stages <-chan amqp.Delivery, queue queues.MessageQueue, onSuccess func(types.Message, queues.MessageQueue)) {
	for d := range stages {
		message, err := unmarshalMessage(d.Body)
		if err != nil {
			panic(err)
		}
		onSuccess(*message, queue)
	}
}

func unmarshalMessage(body []byte) (*types.Message, error) {
	var message types.Message
	err := json.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (p *pipelineService) Run(pipeline *types.Pipeline) error {

	err := p.queue.AddStageToQueue(types.Message{
		CurrentStage: pipeline.Stages[0],
		Pipeline:     *pipeline,
	})
	if err != nil {
		return err
	}

	return nil
}
