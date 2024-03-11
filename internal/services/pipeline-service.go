package services

import (
	"github.com/PiotrFerenc/mash2/internal/queues"
)

type PipelineService interface {
	Run()
}

type pipelineService struct {
	queue queues.MessageQueue
}

func CreatePipelineService(queue queues.MessageQueue) PipelineService {
	return &pipelineService{
		queue: queue,
	}
}

func (p *pipelineService) Run() {

}
