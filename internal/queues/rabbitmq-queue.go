package queues

type queue struct {
}

func CreateRabbitMqMessageQueue() MessageQueue {
	return &queue{}
}

func (queue *queue) Connect() {
	//TODO: implement
}
func (queue *queue) Receive() {
	//TODO: implement
}
func (queue *queue) Publish() {
	//TODO: implement
}
