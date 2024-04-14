package types

import (
	apitypes "github.com/PiotrFerenc/mash2/api/types"
	"github.com/google/uuid"
)

func MapForeachBody(body apitypes.ForeachBody) ForeachBody {
	process := ForeachBody{
		Id:          uuid.New(),
		Steps:       MapToForeachSteps(body.Stages),
		Error:       "",
		CurrentStep: ForeachStep{},
		Parameters:  body.Parameters,
		Status:      Waiting,
	}

	return process
}
func MapToForeachSteps(steps []apitypes.ForeachStage) []ForeachStep {
	foreachSteps := make([]ForeachStep, len(steps))
	for i, step := range steps {
		foreachSteps[i] = MapToForeachStep(step)
	}
	return foreachSteps

}
func MapToForeachStep(step apitypes.ForeachStage) ForeachStep {
	return ForeachStep{
		Id:       uuid.New(),
		Sequence: step.Sequence,
		Action:   step.Action,
		Name:     step.Name,
		Status:   Waiting,
	}
}

func MapToStep(forEachStep ForeachStep) *Step {
	return &Step{
		Id:       forEachStep.Id,
		Sequence: forEachStep.Sequence,
		Action:   forEachStep.Action,
		Name:     forEachStep.Name,
		Status:   forEachStep.Status,
	}
}

type ForeachBody struct {
	Id          uuid.UUID
	Steps       []ForeachStep
	Error       string
	CurrentStep ForeachStep
	Parameters  map[string]interface{}
	Status      StepStatus
}

type ForeachStep struct {
	Id       uuid.UUID
	Sequence int
	Action   string
	Name     string
	Status   StepStatus
}
