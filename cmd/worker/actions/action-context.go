package actions

import "errors"

type ActionContext struct {
	Parameters map[string]string `json:"parameters"`
}

func (action *ActionContext) GetProperty(name string) (string, error) {
	value, ok := action.Parameters[name]
	if !ok {
		return " ", errors.New("key not found")
	}

	return value, nil
}
