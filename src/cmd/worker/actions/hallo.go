package actions

import (
	"errors"
	"github.com/PiotrFerenc/mash2/src/api/types"
	"log"
)

type hallo struct {
}

func CreateHalloAction() Action {
	return &hallo{}
}

func (action *hallo) Inputs() []Property {
	output := make([]Property, 1)
	output[0] = Property{
		Name: "name",
		Type: "text",
	}
	return output
}

func (action *hallo) Outputs() []Property {
	output := make([]Property, 1)
	output[0] = Property{
		Name: "hallo",
		Type: "text",
	}
	return output
}

func (action *hallo) Execute(message types.Message) (types.Message, error) {
	name, err := message.GetString("name")
	if err != nil {
		return types.Message{}, errors.New("name property not found.")
	}

	msg := "Hallo " + name
	log.Print(msg)
	_, err = message.SetString("hallo", msg)

	return message, nil
}
