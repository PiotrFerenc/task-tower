package types

import "github.com/google/uuid"

type Input struct {
	Id          uuid.UUID
	Name        string
	Type        string
	Description string
	Validation  string
	Value       string
}
