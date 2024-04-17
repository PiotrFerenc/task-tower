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
			Type:        actions.Text,
			Description: "Text to display",
			Validation:  "required",
		},
	}
}

func (action *console) GetCategoryName() string {
	return "console"
}
func (action *console) Inputs() []actions.Property {
	return []actions.Property{
		action.text,
	}
}

func (action *console) Outputs() []actions.Property {
	return []actions.Property{}
}

// Execute executes the action by getting a string value from the `text` property of the `pipeline` parameter,
// logging the value, and returning the `pipeline` and `nil` error if successful.
// If an error occurs while retrieving the string value, it returns an empty `types.Process` and the error.
//
// Params:
//
//	pipeline: The `types.Process` object representing the current process.
//
// Returns:
//
//	types.Process: The modified `pipeline` object.
//	error: The error, if any, occurred during execution.
func (action *console) Execute(pipeline types.Process) (types.Process, error) {
	text, err := action.text.GetStringFrom(&pipeline)
	if err != nil {
		return types.Process{}, err
	}
	log.Print(text)

	return pipeline, nil
}
