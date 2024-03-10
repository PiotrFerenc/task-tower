package types

type Stage struct {
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"`
}
