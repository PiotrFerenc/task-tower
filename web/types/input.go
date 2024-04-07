package types

import "github.com/google/uuid"

type Input struct {
	Id          uuid.UUID
	Name        string
	DisplayName string
	Type        string
	Description string
	Validation  string
	Value       string
}
