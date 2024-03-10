package types

type Step struct {
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"`
}
