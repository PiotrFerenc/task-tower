package actions

import (
	"github.com/PiotrFerenc/mash2/internal/types"
)

type Action interface {
	Inputs() []Property
	Outputs() []Property
	Execute(process types.Pipeline) (types.Pipeline, error)
}

type Property struct {
	Name string
	Type string
}
