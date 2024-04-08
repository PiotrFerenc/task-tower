package repositories

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/web/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StepsRepository interface {
	GetSteps(processId uuid.UUID) ([]types.Step, error)
	Save(actionName string, pipelineId uuid.UUID) (uuid.UUID, error)
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
func (repo *stepRepository) Save(actionName string, pipelineId uuid.UUID) (uuid.UUID, error) {
	var maxSequence int64
	repo.Database.Model(&types.Step{}).Where("value > ?", 10).Count(&maxSequence)
	action := actionName
	name := fmt.Sprintf("%s_%d", actionName, maxSequence)

	step := &types.Step{
		ID:         uuid.New(),
		Sequence:   maxSequence + 1,
		Action:     action,
		Name:       name,
		PipelineID: pipelineId,
	}
	repo.Database.Create(step)

	return step.ID, nil
}
