package services

import (
	"github.com/PiotrFerenc/mash2/internal/repositories"
	"github.com/PiotrFerenc/mash2/internal/types"
)

type ProcessService interface {
	MarkAsStarted(types *types.Process)
	TaskFinished(types *types.Process)
	MarkAsDone(types *types.Process)
	MarkAsFailed(types *types.Process, err string)
}
type processService struct {
	repository repositories.ProcessRepository
}

// CreateProcessService creates a new instance of ProcessService with the given repository.
func CreateProcessService(repository repositories.ProcessRepository) ProcessService {
	return &processService{
		repository: repository,
	}
}

// MarkAsStarted saves the pipeline in the repository.
func (process *processService) MarkAsStarted(pipeline *types.Process) {
	process.repository.Save(*pipeline)
}

// TaskFinished sets the status of the pipeline to "Processing" and its current step to "Done".
// It then updates the status of the pipeline in the repository.
func (process *processService) TaskFinished(pipeline *types.Process) {
	pipeline.Status = types.Processing
	pipeline.CurrentStep.Status = types.Done
	process.repository.UpdateStatus(*pipeline)
}

// MarkAsDone sets the status of the pipeline and its current step to "Done".
// It then updates the status of the pipeline in the repository.
func (process *processService) MarkAsDone(pipeline *types.Process) {
	pipeline.Status = types.Done
	pipeline.CurrentStep.Status = types.Done
	process.repository.UpdateStatus(*pipeline)
}

// MarkAsFailed sets the status of the pipeline and its current step to "Fail",
// and updates the error message. It then updates the status of the pipeline in the repository.
// func (process *processService) MarkAsFailed(pipeline *types.Process, err string) {}
func (process *processService) MarkAsFailed(pipeline *types.Process, err string) {
	pipeline.Status = types.Fail
	pipeline.Error = err
	pipeline.CurrentStep.Status = types.Fail
	process.repository.UpdateStatus(*pipeline)

}
