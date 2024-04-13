package types

import (
	"fmt"
	apitypes "github.com/PiotrFerenc/mash2/api/types"
	"github.com/gobeam/stringy"
	"github.com/google/uuid"
	"sort"
	"strconv"
)

type Pipeline struct {
	Id          uuid.UUID
	Steps       []Step
	Error       string
	CurrentStep Step
	Parameters  map[string]interface{}
	Status      StepStatus
}

type Step struct {
	Id          uuid.UUID
	Sequence    int
	Action      string
	Name        string
	Status      StepStatus
	SubPipeline *Pipeline
}

type StepStatus int

const (
	Waiting StepStatus = iota
	Processing
	Done
	Fail
)

func (message *Pipeline) GetInternalName(propertyName string) string {
	str := stringy.New(message.CurrentStep.Name)
	internalName := str.CamelCase("?", "")
	internalName = stringy.New(internalName).ToLower()
	return fmt.Sprintf("%s.%s", internalName, propertyName)
}

func (message *Pipeline) NewFolder(path string) string {
	return fmt.Sprintf("%s/%s", path, uuid.NewString())
}

func (message *Pipeline) SetInt(name string, value int) {
	message.Parameters[message.GetInternalName(name)] = strconv.Itoa(value)
}
func (message *Pipeline) SetString(name, value string) {
	message.Parameters[message.GetInternalName(name)] = value
}
func NewProcessFromPipeline(pipeline *apitypes.Pipeline) *Pipeline {
	process := mapToPipeline(pipeline)

	return process
}

func mapToPipeline(pipeline *apitypes.Pipeline) *Pipeline {

	process := &Pipeline{
		Id:         uuid.New(),
		Parameters: pipeline.Parameters,
		Steps:      make([]Step, len(pipeline.Stages)),
		Status:     Processing,
	}

	sort.SliceStable(pipeline.Stages, func(i, j int) bool {
		return pipeline.Stages[i].Sequence < pipeline.Stages[j].Sequence
	})

	for i, stage := range pipeline.Stages {
		process.Steps[i] = Step{
			Id:       uuid.New(),
			Sequence: stage.Sequence,
			Action:   stage.Action,
			Name:     stage.Name,
			Status:   Waiting,
		}
		if stage.SubPipeline != nil {
			process.Steps[i].SubPipeline = mapToPipeline(stage.SubPipeline)
		}
	}

	if len(process.Steps) > 0 {
		process.CurrentStep = process.Steps[0]
		process.CurrentStep.Status = Processing
	}
	return process
}
