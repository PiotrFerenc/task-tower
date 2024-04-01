package actions

import (
	"github.com/PiotrFerenc/mash2/internal/types"
	"log"
)

type console struct {
}

func CreateConsoleAction() Action {
	return &console{}
}

func (action *console) Inputs() []Property {
	return []Property{
		{
			Name: "text",
			Type: "text",
		}}
}

func (action *console) Outputs() []Property {
	return []Property{
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
