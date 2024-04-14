package types

import "github.com/go-playground/validator/v10"

type Pipeline struct {
	Stages     []Stage                `json:"stages" validate:"required"`
	Parameters map[string]interface{} `json:"parameters" validate:"required"`
}

func (p *Pipeline) Validate() error {
	validate := validator.New()

	if err := validate.Struct(p); err != nil {
		return err
	}

	for _, stage := range p.Stages {
		if err := validate.Struct(stage); err != nil {
			return err
		}
	}
	return nil
}
