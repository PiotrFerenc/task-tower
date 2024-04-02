package repositories

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var (
	processes []*types.Pipeline
)

type repository struct {
	Database *gorm.DB
}

func CreatePostgresRepository(config *configuration.DatabaseConfig) ProcessRepository {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)

	db, err := gorm.Open(driver.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&ProcessEntity{}, &StepEntity{})
	if err != nil {
		log.Fatal(err)
	}
	return &repository{Database: db}
}

func (r *repository) UpdateStatus(pipeline types.Pipeline) {
	p := ProcessEntity{
		ID: pipeline.Id,
	}
	r.Database.Model(&p).Update("Status", pipeline.Status)
	s := StepEntity{
		ID:        pipeline.CurrentStep.Id,
		ProcessId: p.ID,
	}
	r.Database.Model(&s).Update("Status", pipeline.CurrentStep.Status)
}
func (r *repository) Save(pipeline types.Pipeline) {
	p := ProcessEntity{
		ID:     pipeline.Id,
		Status: int(pipeline.Status),
	}
	r.Database.Save(&p)
	for _, step := range pipeline.Steps {
		s := StepEntity{
			ID:        step.Id,
			ProcessId: pipeline.Id,
			Status:    int(step.Status),
			Name:      step.Name,
		}
		r.Database.Save(&s)
	}
}

type ProcessEntity struct {
	ID     uuid.UUID `gorm:"primaryKey"`
	Status int
}
type StepEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	ProcessId uuid.UUID
	Status    int
	Name      string
}
