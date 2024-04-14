package types

type Stage struct {
	Sequence    int          `json:"sequence" validate:"required,gte=0"`
	Action      string       `json:"action" validate:"required"`
	Name        string       `json:"name" validate:"required"`
	SubPipeline *ForeachBody `json:"body"`
}
type ForeachStage struct {
	Sequence int    `json:"sequence" validate:"required,gte=0"`
	Action   string `json:"action" validate:"required"`
	Name     string `json:"name" validate:"required"`
}
type ForeachBody struct {
	Stages     []ForeachStage         `json:"stages" validate:"required"`
	Parameters map[string]interface{} `json:"parameters" validate:"required"`
}
