package repositories

import (
	"context"
	"fmt"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type mongoDbRepository struct {
	collection *mongo.Collection
}

const collectionName = "pipelines"

// CreateMongoDbRepository initializes a new MongoDB repository for process entities.
func CreateMongoDbRepository(config *configuration.DatabaseConfig) ProcessRepository {
	connectionString := fmt.Sprintf("mongodb://%s:%s", config.DbHost, config.DbPort)
	clientOptions := options.Client().ApplyURI(connectionString)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &mongoDbRepository{
		client.Database(config.DbName).Collection(collectionName),
	}
}

// UpdateStatus updates the status of a process pipeline in the MongoDB repository.
//
// Parameters:
//
//	pipeline: The process pipeline object containing the updated status.
func (m *mongoDbRepository) UpdateStatus(pipeline types.Process) {
	filter := bson.D{{"_id", pipeline.Id}}

	update := bson.D{
		{"$set", bson.D{
			{"Status", pipeline.Status},
		}},
	}
	result, err := m.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	if result.ModifiedCount == 0 {
		log.Fatal(result.UpsertedID)
	}
}

// Save saves a process pipeline in the MongoDB repository.
//
// Parameters:
//
//	pipeline: The process pipeline object to be saved.
func (m *mongoDbRepository) Save(pipeline types.Process) {
	ctx := context.Background()
	_, err := m.collection.InsertOne(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
}

// GetById retrieves a process pipeline from the MongoDB repository based on the given processId.
// Parameters:
//
//	processId: The unique identifier of the process pipeline to retrieve.
//
// Returns:
//
//	ProcessEntity: The retrieved process pipeline entity.
//	error: An error if the retrieval fails.
func (m *mongoDbRepository) GetById(processId uuid.UUID) (ProcessEntity, error) {

	ctx := context.Background()
	result := ProcessEntity{}
	err := m.collection.FindOne(ctx, bson.M{"_id": processId.String()}).Decode(result)
	if err != nil {
		return ProcessEntity{}, err
	}
	return result, nil
}
