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
	UpdateParameters(parameters map[string]interface{}) error
	AddParameters(stepId uuid.UUID, parameters map[string]interface{}) error
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

func (repo *parametersRepository) UpdateParameters(parameters map[string]interface{}) error {
	for key, value := range parameters {
		repo.Database.Model(&types.Parameters{}).Where("ID = ?", key).Update("value", value)
	}
	return nil
}
func (repo *parametersRepository) AddParameters(stepId uuid.UUID, parameters map[string]interface{}) error {
	for key, value := range parameters {
		repo.Database.Model(&types.Parameters{}).Create(types.Parameters{
			ID:     uuid.New(),
			Key:    key,
			Value:  value.(string),
			StepID: stepId,
		})
	}
	return nil
}
