package actions

import "log"

type Hallo struct {
}

func (receiver Hallo) Execute(parameters ActionContext) string {
	name := parameters.Parameters["name"]

	msg := "Hallo " + name
	log.Print(msg)
	return msg
}
