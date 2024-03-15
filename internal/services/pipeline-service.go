package services

import (
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/queues"
)

type PipelineService interface {
	Run(pipeline *types.Pipeline) error
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

func start(pipeline *types.Pipeline, start func(pipeline *types.Pipeline)) {
	start(pipeline)
}

func onSuccess(os func()) {

}
func onFail(fail func()) {

}

func (p *pipelineService) Run(pipeline *types.Pipeline) error {

	err := p.queue.Publish(pipeline.Stages[0])
	if err != nil {
		return err
	}
	// saga provider
	// dodac kroki do bazy
	// pobrac pierwszy
	// wyslac do kolejki -> Execute_action
	//start(pipeline, func(pipeline *types.Pipeline) {
	//	err := p.queue.Publish(types.Stage{})
	//	if err != nil {
	//		return
	//	}
	//	p.processService.MarkAsStarted()
	//})
	//
	//onSuccess(func() {
	//	p.queue.Subscribe()
	//})
	//
	//onFail(func() {
	//	p.queue.Subscribe()
	//	p.processService.MarkAsFailed()
	//})

	return nil
}
