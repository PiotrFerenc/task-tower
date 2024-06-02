package repositories

import "github.com/PiotrFerenc/mash2/internal/events"

type EventStore interface {
	AddEvent(event events.Event) error
}
