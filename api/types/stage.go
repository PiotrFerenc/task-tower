package types

type Stage struct {
	Sequence    int          `json:"sequence" validate:"required,gte=0"`
	Action      string       `json:"action" validate:"required"`
	Name        string       `json:"name" validate:"required"`
	SubPipeline *ForeachBody `json:"pipeline"`
}
type ForeachStage struct {
	Sequence int    `json:"sequence" validate:"required,gte=0"`
	Action   string `json:"action" validate:"required"`
	Name     string `json:"name" validate:"required"`
}
