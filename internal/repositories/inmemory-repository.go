package repositories

import (
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
)

var (
	processes []*types.Pipeline
)

type repository struct {
}

func CreateInMemoryRepository() ProcessRepository {
	processes = make([]*types.Pipeline, 1)
	return &repository{}
}

func (r *repository) GetProcess(id uuid.UUID) *types.Pipeline {
	for _, process := range processes {
		if process.Id == id {
			return process
		}
	}
	return nil
}

func (r *repository) UpdateStatus(pipeline *types.Pipeline) {
	for i, process := range processes {
		if process.Id == pipeline.Id {
			processes[i] = pipeline
			return
		}
	}
}
func (r *repository) Start(pipeline *types.Pipeline) {
	processes[0] = pipeline
}
