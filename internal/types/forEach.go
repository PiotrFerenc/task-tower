package types

import (
	apitypes "github.com/PiotrFerenc/mash2/api/types"
	"github.com/google/uuid"
)

// MapForeachBody maps a ForeachBody object to a ForeachBody object with additional fields.
func MapForeachBody(body apitypes.ForeachBody) ForeachBody {
	process := ForeachBody{
		Id:          uuid.New(),
		Steps:       MapToForeachSteps(body.Tasks),
		Error:       "",
		CurrentStep: ForeachStep{},
		Parameters:  body.Parameters,
		Status:      Waiting,
	}

	return process
}

// MapToForeachSteps maps a slice of ForeachTasks to a slice of ForeachSteps.
//
// Parameters:
//
//	steps: The slice of ForeachTasks to be mapped.
//
// Returns:
//
//	[]ForeachStep: The mapped slice of ForeachSteps.
func MapToForeachSteps(steps []apitypes.ForeachTask) []ForeachStep {
	foreachSteps := make([]ForeachStep, len(steps))
	for i, step := range steps {
		foreachSteps[i] = MapToForeachStep(step)
	}
	return foreachSteps

}

// MapToForeachStep maps a ForeachTask to a ForeachStep.
//
// Parameters:
//
//	step: The ForeachTask to be mapped.
//
// Returns:
//
//	ForeachStep: The mapped ForeachStep.
func MapToForeachStep(step apitypes.ForeachTask) ForeachStep {
	return ForeachStep{
		Id:       uuid.New(),
		Sequence: step.Sequence,
		Action:   step.Action,
		Name:     step.Name,
		Status:   Waiting,
	}
}

// MapToStep maps a ForeachStep to a Step.
//
// Parameters:
//
//	forEachStep: The ForeachStep to be mapped.
//
// Returns:
//
//	*Step: A pointer to the mapped Step.
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
