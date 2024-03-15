package queues

import "github.com/PiotrFerenc/mash2/api/types"

type MessageQueue interface {
	Connect() error
	Publish(message types.Stage) error
	Subscribe()
	CreateQueue(name string) error
}
