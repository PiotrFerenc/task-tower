package types

import (
	"errors"
	"fmt"
	"github.com/gobeam/stringy"
	"github.com/google/uuid"
	"github.com/valyala/fasttemplate"
	"strconv"
)

type Process struct {
	Id          uuid.UUID
	Steps       []Step
	Error       string
	CurrentStep Step
	Parameters  map[string]interface{}
}

type Step struct {
	Id       uuid.UUID
	Sequence int
	Action   string
	Name     string
}

func (message *Process) GetInternalName(propertyName string) string {
	str := stringy.New(message.CurrentStep.Name)
	internalName := str.CamelCase("?", "")
	internalName = stringy.New(internalName).ToLower()
	return fmt.Sprintf("%s.%s", internalName, propertyName)
}

func (message *Process) GetString(name string) (string, error) {
	internalName := message.GetInternalName(name)
	parameter, ok := message.Parameters[internalName]
	if !ok {
		return " ", errors.New("key not found")
	}
	value := parameter.(string)

	template := fasttemplate.New(value, "{{", "}}")
	value = template.ExecuteString(message.Parameters)
	message.Parameters[internalName] = value
	return value, nil
}

func (message *Process) NewFolder(path string) string {
	return fmt.Sprintf("%s/%s", path, uuid.NewString())
}

func (message *Process) GetInt(name string) (int, error) {
	value, err := message.GetString(name)
	if err != nil {
		return 0, err
	}
	conv, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return conv, nil
}
func (message *Process) SetInt(name string, value int) (*Process, error) {
	message.Parameters[message.GetInternalName(name)] = strconv.Itoa(value)
	return message, nil
}
func (message *Process) SetString(name, value string) (*Process, error) {
	message.Parameters[message.GetInternalName(name)] = value
	return message, nil
}
