package types

type Stage struct {
	Sequence int    `json:"sequence"`
	Action   string `json:"action"`
	Name     string `json:"name"`
}
