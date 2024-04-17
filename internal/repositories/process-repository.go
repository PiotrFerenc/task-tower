package repositories

import (
	types "github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
)

type ProcessRepository interface {
	UpdateStatus(pipeline types.Process)
	Save(pipeline types.Process)
	GetById(processId uuid.UUID) (ProcessEntity, error)
}
