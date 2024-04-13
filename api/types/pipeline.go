package types

import "github.com/go-playground/validator/v10"

type ForeachBody struct {
	Stages     []ForeachStage         `json:"stages" validate:"required"`
	Parameters map[string]interface{} `json:"parameters" validate:"required"`
}
type Pipeline struct {
	Stages     []Stage                `json:"stages" validate:"required"`
	Parameters map[string]interface{} `json:"parameters" validate:"required"`
}

func (p *Pipeline) Validate() error {
	validate := validator.New()
	validate.RegisterStructValidation(StageStructLevelValidation, Stage{})

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
func StageStructLevelValidation(sl validator.StructLevel) {
	stage := sl.Current().Interface().(Stage)

	if stage.Action == "for-each" && stage.SubPipeline == nil {
		sl.ReportError(stage.SubPipeline, "SubPipeline", "SubPipeline", "required", "")
	}
}
