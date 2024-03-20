package types

type Stage struct {
	Order      int               `json:"order"`
	Action     string            `json:"action"`
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"`
}
