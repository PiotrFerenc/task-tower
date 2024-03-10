package executor

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
)

type MapExecutor struct {
}

type executor struct {
}

func CreateMapExecutor() Executor {
	return &executor{}
}

func (executor *executor) Execute(actionName string, parameters actions.ActionContext) error {
	a := map[string]actions.Action{
		"hallo": actions.Hallo{},
		"sleep": actions.Sleep{},
	}
	action, exist := a[actionName]
	if !exist {
		return fmt.Errorf("action %v not found", actionName)
	}
	action.Execute(parameters)
	return nil
}
