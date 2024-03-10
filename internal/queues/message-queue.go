package queues

type MessageQueue interface {
	Connect()
	Publish()
	Receive()
}
