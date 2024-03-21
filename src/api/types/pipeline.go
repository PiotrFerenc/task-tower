package types

type Pipeline struct {
	Stages     []Stage                `json:"stages"`
	Parameters map[string]interface{} `json:"parameters"`
}
