package others

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/types"
	"log"
)

type console struct {
}

func CreateConsoleAction() actions.Action {
	return &console{}
}

func (action *console) Inputs() []actions.Property {
	return []actions.Property{
		{
			Name: "text",
			Type: "text",
		}}
}

func (action *console) Outputs() []actions.Property {
	return []actions.Property{
		{
			Name: "text",
			Type: "text",
		},
	}
}

func (action *console) Execute(message types.Pipeline) (types.Pipeline, error) {
	name, _ := message.GetString("text")

	log.Print(name)
	_, _ = message.SetString("text", name)

	return message, nil
}
