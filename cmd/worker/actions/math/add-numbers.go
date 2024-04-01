package math

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/types"
)

type addnumbers struct {
	a actions.Property
	b actions.Property
	c actions.Property
}

// CreateAddNumbers This is a function that initializes an instance of the addnumbers struct.
// It returns a pointer to the addnumbers instance.
// This is useful when we don't want to pass the struct by value in subsequent calls.
func CreateAddNumbers() actions.Action {
	return &addnumbers{
		actions.Property{
			Name:        "a",
			Type:        "number",
			Description: "a",
			Validation:  "required,number",
		},
		actions.Property{
			Name:        "b",
			Type:        "number",
			Description: "b",
			Validation:  "required,number",
		},
		actions.Property{
			Name:        "c",
			Type:        "number",
			Description: "c",
			Validation:  "",
		},
	}
}

// Inputs The Inputs() method returns a slice of Property structure.
// The Property structure includes two fields: Name and Type, both of which are strings.
// These property structures are created for two inputs, 'a' and 'b', of 'number' type.
// It then returns these properties.
func (action *addnumbers) Inputs() []actions.Property {
	return []actions.Property{
		action.a, action.b,
	}

}

// Outputs The Outputs() method returns a slice of Property structure.
// It creates a property structure for an output, 'c', of 'number' type and returns it.
func (action *addnumbers) Outputs() []actions.Property {
	return []actions.Property{action.c}
}

// Execute The Execute() method receives a parameter of types.Message type and returns (types.Message, error).
func (action *addnumbers) Execute(pipeline types.Pipeline) (types.Pipeline, error) {

	// In this Execute() method, first, it retrieves the integer values 'a' and 'b' from the pipeline.
	a, err := action.a.GetIntFrom(&pipeline)
	if err != nil {
		return pipeline, err
	}
	b, err := action.b.GetIntFrom(&pipeline)
	if err != nil {
		return pipeline, err
	}
	//Then it adds them together and sets the resulting 'c' back into the pipeline.
	c := a + b

	//After performing these operations, it returns the updated pipeline and nil for the error value.
	pipeline.SetInt(action.c.Name, c)

	return pipeline, nil

}
