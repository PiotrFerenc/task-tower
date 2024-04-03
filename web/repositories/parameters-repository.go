package repositories

import (
	"github.com/PiotrFerenc/mash2/web/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type parametersRepository struct {
	Database *gorm.DB
}

type ParametersRepository interface {
	GetParameters(stepId uuid.UUID) []types.Parameters
}

func CreateParametersRepository(connection *gorm.DB) ParametersRepository {
	return &parametersRepository{
		Database: connection,
	}
}

func (repo *parametersRepository) GetParameters(stepId uuid.UUID) []types.Parameters {
	var parameters []types.Parameters
	repo.Database.Where("step_id = ?", stepId).Find(&parameters)
	return parameters
}
