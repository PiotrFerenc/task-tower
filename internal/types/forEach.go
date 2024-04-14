package types

import (
	apitypes "github.com/PiotrFerenc/mash2/api/types"
	"github.com/google/uuid"
)

func MapForeachBody(body *apitypes.ForeachBody) *ForeachBody {
	process := &ForeachBody{
		Id:          uuid.New(),
		Steps:       make([]ForeachStep, 0),
		Error:       "",
		CurrentStep: ForeachStep{},
		Parameters:  body.Parameters,
		Status:      Waiting,
	}

	return process
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
