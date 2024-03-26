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
	output := make([]Property, 1)
	output[0] = Property{
		Name: "text",
		Type: "text",
	}
	return output
}

func (action *console) Outputs() []Property {
	output := make([]Property, 1)
	output[0] = Property{
		Name: "text",
		Type: "text",
	}
	return output
}

func (action *console) Execute(message types.Process) (types.Process, error) {
	name, _ := message.GetString("text")

	log.Print(name)
	_, _ = message.SetString("text", name)

	return message, nil
}
