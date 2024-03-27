package repositories

import (
	types "github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
)

type ProcessRepository interface {
	UpdateStatus(pipeline *types.Pipeline)
	Start(pipeline *types.Pipeline)
	GetProcess(id uuid.UUID) *types.Pipeline
}
