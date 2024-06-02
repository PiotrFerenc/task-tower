package events

import (
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/google/uuid"
	"time"
)

type Event struct {
	CreatedAt time.Time
	StreamId  uuid.UUID
}

type ProcessEvent interface {
	ProcessEvent() *Event
}

type ProcessCreated struct {
	Process types.Process
}

func (p ProcessCreated) ProcessEvent() *Event {
	return &Event{
		CreatedAt: time.Now(),
		StreamId:  p.Process.Id,
	}
}
