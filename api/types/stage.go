package types

type Task struct {
	Sequence    int          `json:"sequence" validate:"required,gte=0"`
	Action      string       `json:"action" validate:"required"`
	Name        string       `json:"name" validate:"required"`
	SubPipeline *ForeachBody `json:"body"`
}
type ForeachTask struct {
	Sequence int    `json:"sequence" validate:"required,gte=0"`
	Action   string `json:"action" validate:"required"`
	Name     string `json:"name" validate:"required"`
}
type ForeachBody struct {
	Tasks      []ForeachTask          `json:"Tasks" validate:"required"`
	Parameters map[string]interface{} `json:"parameters" validate:"required"`
}
