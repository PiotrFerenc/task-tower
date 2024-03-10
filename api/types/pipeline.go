package types

type Pipeline struct {
	Steps []Stage `json:"stages"`
}
