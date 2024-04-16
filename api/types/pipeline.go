package types

import "github.com/go-playground/validator/v10"

type Pipeline struct {
	Tasks      []Task                 `json:"Tasks" validate:"required"`
	Parameters map[string]interface{} `json:"parameters" validate:"required"`
}

func (p *Pipeline) Validate() error {
	validate := validator.New()

	if err := validate.Struct(p); err != nil {
		return err
	}

	for _, Task := range p.Tasks {
		if err := validate.Struct(Task); err != nil {
			return err
		}
	}
	return nil
}
