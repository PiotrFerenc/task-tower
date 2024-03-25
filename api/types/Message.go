package types

type Message struct {
	Error        string
	CurrentStage Stage
	Pipeline     Pipeline
}
