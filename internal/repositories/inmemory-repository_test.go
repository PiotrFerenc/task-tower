package repositories

import (
	"testing"

	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetProcess(t *testing.T) {
	tests := []struct {
		name      string
		processes []*types.Pipeline
		id        uuid.UUID
		expected  *types.Pipeline
	}{
		{
			name:      "find existing process",
			processes: []*types.Pipeline{{Id: uuid.MustParse("ba8e1fa4-8776-4b71-a553-d416232e8571")}},
			id:        uuid.MustParse("ba8e1fa4-8776-4b71-a553-d416232e8571"),
			expected:  &types.Pipeline{Id: uuid.MustParse("ba8e1fa4-8776-4b71-a553-d416232e8571")},
		},
		{
			name:      "process not found",
			processes: []*types.Pipeline{{Id: uuid.MustParse("ba8e1fa4-8776-4b71-a553-d416232e8571")}},
			id:        uuid.MustParse("5779be80-85dd-4e71-a598-44d9ef7419b5"),
			expected:  nil,
		},
		{
			name:      "empty repository",
			processes: []*types.Pipeline{},
			id:        uuid.MustParse("ba8e1fa4-8776-4b71-a553-d416232e8571"),
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processes = tt.processes
			r := &repository{}
			result := r.GetProcess(tt.id)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	items := []*types.Pipeline{
		{Id: uuid.New(), Status: types.Done},
		{Id: uuid.New(), Status: types.Fail},
	}

	tests := []struct {
		name              string
		initialpipelines  []*types.Pipeline
		updatePipeline    *types.Pipeline
		expectedPipelines []*types.Pipeline
	}{
		{
			name:             "Successful Update",
			initialpipelines: items,
			updatePipeline:   &types.Pipeline{Id: items[0].Id, Status: types.Waiting},
			expectedPipelines: []*types.Pipeline{
				{Id: items[0].Id, Status: types.Waiting},
				{Id: items[1].Id, Status: types.Fail},
			},
		},
		{
			name:              "Update Non-Existing Pipeline",
			initialpipelines:  items,
			updatePipeline:    &types.Pipeline{Id: uuid.New(), Status: types.Fail},
			expectedPipelines: items,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processes = tt.initialpipelines
			r := &repository{}
			r.UpdateStatus(tt.updatePipeline)

			for i, pipeline := range processes {
				if pipeline.Id != tt.expectedPipelines[i].Id ||
					pipeline.Status != tt.expectedPipelines[i].Status {
					t.Errorf("UpdateStatus() = %v, want %v", pipeline, tt.expectedPipelines[i])
				}
			}
		})
	}
}

func TestSave(t *testing.T) {
	repo := CreateInMemoryRepository()

	t.Run("single save", func(t *testing.T) {
		pipeline := &types.Pipeline{
			Id: uuid.New(),
		}
		repo.Save(pipeline)

		savedPipeline := repo.GetProcess(pipeline.Id)
		assert.Equal(t, pipeline.Id, savedPipeline.Id)
	})

	t.Run("multiple saves", func(t *testing.T) {
		firstPipeline := &types.Pipeline{
			Id: uuid.New(),
		}
		secondPipeline := &types.Pipeline{
			Id: uuid.New(),
		}

		repo.Save(firstPipeline)
		repo.Save(secondPipeline)

		savedFirstPipeline := repo.GetProcess(firstPipeline.Id)
		assert.Equal(t, firstPipeline.Id, savedFirstPipeline.Id)

		savedSecondPipeline := repo.GetProcess(secondPipeline.Id)
		assert.Equal(t, secondPipeline.Id, savedSecondPipeline.Id)
	})

	t.Run("save and update", func(t *testing.T) {
		pipeline := &types.Pipeline{
			Id: uuid.New(),
		}
		repo.Save(pipeline)

		pipeline.Status = types.Processing
		repo.UpdateStatus(pipeline)

		savedPipeline := repo.GetProcess(pipeline.Id)

		assert.Equal(t, pipeline.Id, savedPipeline.Id)
		assert.Equal(t, pipeline.Status, savedPipeline.Status)
	})
}
