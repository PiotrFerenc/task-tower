package actions

import (
	"github.com/PiotrFerenc/mash2/api/types"
)

type Action interface {
	Inputs() []Property
	Outputs() []Property
	Execute(message types.Message) (types.Message, error)
}

type Property struct {
	Name string
	Type string
}
