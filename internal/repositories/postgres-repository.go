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
	processes []*types.Process
)

type repository struct {
	Database *gorm.DB
}

// CreatePostgresRepository takes a DatabaseConfig and returns a ProcessRepository that interacts with a PostgreSQL database.
// The function establishes a connection to the database using the provided configuration, performs necessary migrations,
// and returns a repository implementation that uses the connection to save and update process and step entities.
// In case of any error during connection or migrations, the function logs the error and terminates the program.
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
func (r *repository) GetById(processId uuid.UUID) (ProcessEntity, error) {
	processes := ProcessEntity{}
	r.Database.Find(&processes, processId)
	return processes, nil
}

// UpdateStatus updates the status of a pipeline in the repository.
// It takes a types.Process object and updates the corresponding ProcessEntity and StepEntity records in the database.
// The function updates the "Status" and "Error" fields of the ProcessEntity with the values from the pipeline object.
// Additionally, it updates the "Status" field of the StepEntity with the status of the current step in the pipeline object.
func (r *repository) UpdateStatus(pipeline types.Process) {
	p := ProcessEntity{
		ID: pipeline.Id,
	}
	r.Database.Model(&p).Updates(map[string]interface{}{
		"Status": pipeline.Status,
		"Error":  pipeline.Error,
	})
	s := StepEntity{
		ID: pipeline.CurrentStep.Id,
	}
	r.Database.Model(&s).Update("Status", pipeline.CurrentStep.Status)
}

// Save saves a pipeline to the repository.
// It takes a types.Process object and saves the corresponding ProcessEntity and StepEntity records in the database.
// The function creates a new ProcessEntity with the ID and Status from the pipeline object, and saves it using the Database.
// Then, it iterates over the Steps in the pipeline object and creates a new StepEntity for each step, with the ID, ProcessId, Status, and Name from the step object.
// Each StepEntity is then saved using the Database.
// Note: The StepEntity records are not directly connected to the ProcessEntity, but they have a foreign key reference to the ProcessEntity's ID.
// Without duplicating the unchanged declarations above, document the following code:
func (r *repository) Save(pipeline types.Process) {
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
	Error  string
}
type StepEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	ProcessId uuid.UUID
	Status    int
	Name      string
}
