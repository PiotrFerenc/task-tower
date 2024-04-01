package repositories

import (
	types "github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
)

type ProcessRepository interface {
	UpdateStatus(pipeline types.Pipeline)
	Save(pipeline types.Pipeline)
	GetProcess(id uuid.UUID) *types.Pipeline
}
