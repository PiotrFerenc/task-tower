package services

import (
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/queues"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type PipelineService interface {
	Run(pipeline *types.Pipeline) error
}

type pipelineService struct {
	queue          queues.MessageQueue
	processService ProcessService
}

func CreatePipelineService(queue queues.MessageQueue, processService ProcessService) PipelineService {
	s := func(d amqp091.Delivery) {
		log.Printf(" success [x] %s", d.Body)
	}

	f := func(d amqp091.Delivery) {
		log.Printf(" fail [x] %s", d.Body)
	}

	go func(onSucces func(amqp091.Delivery), onFail func(amqp091.Delivery)) {
		ss, _ := queue.WaitingForSucceedStage()
		fs, _ := queue.WaitingForFailedStage()

		var forever chan struct{}

		go func() {
			for d := range ss {
				onSucces(d)
			}
		}()

		go func() {
			for d := range fs {
				onFail(d)
			}
		}()

		<-forever

	}(s, f)
	return &pipelineService{
		queue:          queue,
		processService: processService,
	}
}

func start(pipeline *types.Pipeline, start func(pipeline *types.Pipeline)) {
	start(pipeline)
}

func onSuccess(os func()) {

}
func onFail(fail func()) {

}

func (p *pipelineService) Run(pipeline *types.Pipeline) error {

	err := p.queue.AddStageToQueue(pipeline.Stages[0])
	if err != nil {
		return err
	}

	return nil
}
