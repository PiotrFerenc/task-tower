package types

import (
	"github.com/google/uuid"
)

type Pipeline struct {
	ID   uuid.UUID `gorm:"primaryKey"`
	Name string
}

type Parameters struct {
	ID     uuid.UUID `gorm:"primaryKey"`
	Key    string
	Value  string
	StepID uuid.UUID
}

type Step struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	Sequence   int64
	Action     string
	Name       string
	PipelineID uuid.UUID
}
