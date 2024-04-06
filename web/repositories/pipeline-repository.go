package repositories

import (
	"github.com/PiotrFerenc/mash2/web/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PipelineRepository interface {
	GetAll() ([]types.Pipeline, error)
	Save(name string) (types.Pipeline, error)
	GetById(id uuid.UUID) (types.Pipeline, error)
}

func CreatePipelineRepository(connection *gorm.DB) PipelineRepository {
	return &repository{
		Database: connection,
	}
}

type repository struct {
	Database *gorm.DB
}

func (repo *repository) GetById(id uuid.UUID) (types.Pipeline, error) {
	var result types.Pipeline
	if err := repo.Database.Find(&result, id); err.Error != nil {
		return result, err.Error
	}

	return result, nil
}
func (repo *repository) GetAll() ([]types.Pipeline, error) {
	var result []types.Pipeline
	if err := repo.Database.Find(&result); err.Error != nil {
		return result, err.Error
	}

	return result, nil
}

func (repo *repository) Save(name string) (types.Pipeline, error) {

	pipeline := types.Pipeline{
		ID:   uuid.New(),
		Name: name,
	}

	if err := repo.Database.Create(&pipeline); err != nil {
		return pipeline, err.Error
	}
	return pipeline, nil
}
