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
	ForeachBody ForeachBody
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

	process := &Pipeline{
		Id:         uuid.New(),
		Parameters: pipeline.Parameters,
		Steps:      make([]Step, len(pipeline.Tasks)),
		Status:     Processing,
	}

	sort.SliceStable(pipeline.Tasks, func(i, j int) bool {
		return pipeline.Tasks[i].Sequence < pipeline.Tasks[j].Sequence
	})

	for i, Task := range pipeline.Tasks {
		process.Steps[i] = Step{
			Id:       uuid.New(),
			Sequence: Task.Sequence,
			Action:   Task.Action,
			Name:     Task.Name,
			Status:   Waiting,
		}
		if Task.SubPipeline != nil {
			process.Steps[i].ForeachBody = MapForeachBody(*Task.SubPipeline)
		}
	}

	if len(process.Steps) > 0 {
		process.CurrentStep = process.Steps[0]
		process.Steps = process.Steps[1:]
		process.CurrentStep.Status = Processing
	}
	return process
}
