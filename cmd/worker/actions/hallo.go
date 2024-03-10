package actions

import "log"

type Hallo struct {
}

func (receiver Hallo) Execute(parameters ActionContext) string {
	name, err := parameters.GetProperty("name")
	if err == nil {

	}

	msg := "Hallo " + name
	log.Print(msg)
	return msg
}
