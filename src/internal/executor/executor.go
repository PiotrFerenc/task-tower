package executor

import (
	"github.com/PiotrFerenc/mash2/src/api/types"
)

type Executor interface {
	Execute(actionName string, parameters types.Message) error
}
