package services

import (
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/queues"
)

type PipelineService interface {
	Run(pipeline types.Pipeline) error
}

type pipelineService struct {
	queue          queues.MessageQueue
	processService ProcessService
}

func CreatePipelineService(queue queues.MessageQueue, processService ProcessService) PipelineService {
	return &pipelineService{
		queue:          queue,
		processService: processService,
	}
}

func (p *pipelineService) Run(pipeline types.Pipeline) error {
	var stage = pipeline.Stages[0]
	p.queue.Connect()

	err := p.queue.Publish(stage)
	if err == nil {
		return err
	}
	p.processService.MarkAsStarted()

	p.queue.Receive()

	p.processService.MarkAsDone()
	return nil
}
