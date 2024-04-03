package repositories

import (
	"github.com/PiotrFerenc/mash2/web/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StepsRepository interface {
	GetSteps(processId uuid.UUID) ([]types.Step, error)
}
type stepRepository struct {
	Database *gorm.DB
}

func CreateStepsRepository(connection *gorm.DB) StepsRepository {
	return &stepRepository{Database: connection}
}
func (repo *stepRepository) GetSteps(processId uuid.UUID) ([]types.Step, error) {
	var result []types.Step
	if err := repo.Database.Where(&types.Step{PipelineID: processId}).Find(&result); err != nil {
		return result, err.Error
	}
	return result, nil
}
