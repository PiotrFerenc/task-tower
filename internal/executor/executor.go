package executor

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
)

type Executor interface {
	Execute(actionName string, parameters actions.ActionContext) error
}
