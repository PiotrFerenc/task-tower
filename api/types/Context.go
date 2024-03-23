package types

import (
	"errors"
	"fmt"
	"github.com/gobeam/stringy"
	"github.com/valyala/fasttemplate"
	"strconv"
)

func (message *Message) GetInternalName(propertyName string) string {
	str := stringy.New(message.CurrentStage.Name)
	internalName := str.CamelCase("?", "")
	internalName = stringy.New(internalName).ToLower()
	return fmt.Sprintf("%s.%s", internalName, propertyName)
}

func (message *Message) GetString(name string) (string, error) {
	internalName := message.GetInternalName(name)
	parameter, ok := message.Pipeline.Parameters[internalName]
	if !ok {
		return " ", errors.New("key not found")
	}
	value := parameter.(string)

	template := fasttemplate.New(value, "{{", "}}")
	value = template.ExecuteString(message.Pipeline.Parameters)
	message.Pipeline.Parameters[internalName] = value
	return value, nil
}

func (message *Message) GetInt(name string) (int, error) {
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
func (message *Message) SetInt(name string, value int) (*Message, error) {
	message.Pipeline.Parameters[message.GetInternalName(name)] = strconv.Itoa(value)
	return message, nil
}
func (message *Message) SetString(name, value string) (*Message, error) {
	message.Pipeline.Parameters[message.GetInternalName(name)] = value
	return message, nil
}
