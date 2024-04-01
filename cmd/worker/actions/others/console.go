package others

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/types"
	"log"
)

type console struct {
	text actions.Property
}

func CreateConsoleAction() actions.Action {
	return &console{
		text: actions.Property{
			Name:        "text",
			Type:        "text",
			Description: "Text to display",
			Validation:  "required",
		},
	}
}

func (action *console) Inputs() []actions.Property {
	return []actions.Property{
		action.text,
	}
}

func (action *console) Outputs() []actions.Property {
	return []actions.Property{}
}

func (action *console) Execute(pipeline types.Pipeline) (types.Pipeline, error) {
	text, err := action.text.GetStringFrom(&pipeline)
	if err != nil {
		return types.Pipeline{}, err
	}
	log.Print(text)

	return pipeline, nil
}
