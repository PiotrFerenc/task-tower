package actions

import (
	"errors"
	"fmt"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/valyala/fasttemplate"
	"strconv"
)

type Action interface {
	Inputs() []Property
	Outputs() []Property
	GetCategoryName() string
	Execute(process types.Process) (types.Process, error)
}

type Property struct {
	Name        string
	Type        string
	Description string
	DisplayName string
	Validation  string
}

var (
	validate = validator.New()
)

func (property *Property) GetIntFrom(message *types.Process) (int, error) {
	value, err := property.GetStringFrom(message)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}
func (property *Property) GetStringFrom(message *types.Process) (string, error) {
	internalName := message.getInternalName(property.Name)
	parameter, ok := message.Parameters[internalName]
	if !ok {
		msg := fmt.Sprintf("Key %s not found in %s configuration.", internalName, message.CurrentStep.Name)
		return "", errors.New(msg)
	}
	value := parameter.(string)
	if property.Validation != "" {
		if err := validate.Var(value, property.Validation); err != nil {
			return "", err
		}
	}
	template := fasttemplate.New(value, "{{", "}}")
	value = template.ExecuteString(message.Parameters)
	message.Parameters[internalName] = value
	return value, nil
}

const (
	Text     = "text"
	Number   = "number"
	Loop     = "loop"
	Pipeline = "pipeline"
)
