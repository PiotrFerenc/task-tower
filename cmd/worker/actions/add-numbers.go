package actions

import (
	"github.com/PiotrFerenc/mash2/internal/types"
)

type addnumbers struct {
}

// CreateAddNumbers This is a function that initializes an instance of the addnumbers struct.
// It returns a pointer to the addnumbers instance.
// This is useful when we don't want to pass the struct by value in subsequent calls.
func CreateAddNumbers() Action {
	return &addnumbers{}
}

// Inputs The Inputs() method returns a slice of Property structure.
// The Property structure includes two fields: Name and Type, both of which are strings.
// These property structures are created for two inputs, 'a' and 'b', of 'number' type.
// It then returns these properties.
func (action *addnumbers) Inputs() []Property {
	return []Property{
		{
			Name: "a",
			Type: "number",
		},
		{
			Name: "b",
			Type: "number",
		},
	}
}

// Outputs The Outputs() method returns a slice of Property structure.
// It creates a property structure for an output, 'c', of 'number' type and returns it.
func (action *addnumbers) Outputs() []Property {
	return []Property{
		{
			Name: "c",
			Type: "number",
		}}
}

// Execute The Execute() method receives a parameter of types.Message type and returns (types.Message, error).
func (action *addnumbers) Execute(message types.Pipeline) (types.Pipeline, error) {

	// In this Execute() method, first, it retrieves the integer values 'a' and 'b' from the message.
	a, _ := message.GetInt("a")
	b, _ := message.GetInt("b")

	//Then it adds them together and sets the resulting 'c' back into the message.
	c := a + b

	//After performing these operations, it returns the updated message and nil for the error value.
	_, _ = message.SetInt("c", c)

	return message, nil

}
