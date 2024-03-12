package queues

import "github.com/PiotrFerenc/mash2/api/types"

type MessageQueue interface {
	Connect()
	Publish(message types.Stage) error
	Receive()
}
