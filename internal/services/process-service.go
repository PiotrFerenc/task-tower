package services

import (
	"github.com/PiotrFerenc/mash2/internal/repositories"
)

type ProcessService interface {
	MarkAsStarted()
	MarkAsDone()
	MarkAsFailed()
}
type processService struct {
	repository repositories.ProcessRepository
}

func CreateProcessService(repository repositories.ProcessRepository) ProcessService {
	return &processService{
		repository: repository,
	}
}

func (process *processService) MarkAsStarted() {
	process.repository.UpdateStatus()
}
func (process *processService) MarkAsDone() {
	process.repository.UpdateStatus()
}
func (process *processService) MarkAsFailed() {
	process.repository.UpdateStatus()
}
