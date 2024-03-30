package types

import (
	"github.com/google/uuid"
)

type Pipeline struct {
	Id   uuid.UUID `gorm:"primaryKey"`
	Name string
}

type Parameters struct {
	Id         uuid.UUID `gorm:"primaryKey"`
	Key        string
	Value      string
	PipelineId uuid.UUID
}

type Step struct {
	Id         uuid.UUID `gorm:"primaryKey"`
	Sequence   int
	Action     string
	Name       string
	PipelineId uuid.UUID
}
