package services

import (
	"github.com/PiotrFerenc/mash2/internal/repositories"
	"github.com/PiotrFerenc/mash2/internal/types"
)

type ProcessService interface {
	MarkAsStarted(types *types.Pipeline)
	TaskFinished(types *types.Pipeline)
	MarkAsDone(types *types.Pipeline)
	MarkAsFailed(types *types.Pipeline, err string)
}
type processService struct {
	repository repositories.ProcessRepository
}

func CreateProcessService(repository repositories.ProcessRepository) ProcessService {
	return &processService{
		repository: repository,
	}
}

func (process *processService) MarkAsStarted(pipeline *types.Pipeline) {
	process.repository.Save(*pipeline)
}
func (process *processService) TaskFinished(pipeline *types.Pipeline) {
	pipeline.Status = types.Processing
	pipeline.CurrentStep.Status = types.Done
	process.repository.UpdateStatus(*pipeline)
}

func (process *processService) MarkAsDone(pipeline *types.Pipeline) {
	pipeline.Status = types.Done
	pipeline.CurrentStep.Status = types.Done
	process.repository.UpdateStatus(*pipeline)
}
func (process *processService) MarkAsFailed(pipeline *types.Pipeline, err string) {
	pipeline.Status = types.Fail
	pipeline.Error = err
	pipeline.CurrentStep.Status = types.Fail
	process.repository.UpdateStatus(*pipeline)

}
