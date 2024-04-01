package services

import (
	"github.com/PiotrFerenc/mash2/internal/repositories"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
)

type ProcessService interface {
	MarkAsStarted(types *types.Pipeline)
	MarkAsDone(types *types.Pipeline)
	MarkAsFailed(types *types.Pipeline, err string)
	GetProcess(id uuid.UUID) *types.Pipeline
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
func (process *processService) MarkAsDone(pipeline *types.Pipeline) {
	pipeline.Status = types.Done
	process.repository.UpdateStatus(*pipeline)
}
func (process *processService) MarkAsFailed(pipeline *types.Pipeline, err string) {
	pipeline.Status = types.Fail
	pipeline.Error = err
	process.repository.UpdateStatus(*pipeline)
}
func (process *processService) GetProcess(id uuid.UUID) *types.Pipeline {
	return process.repository.GetProcess(id)
}
