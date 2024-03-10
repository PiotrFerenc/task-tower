package queues

type queue struct {
}

func CreateRabbitMqMessageQueue() MessageQueue {
	return &queue{}
}

func (queue *queue) Connect() {}
func (queue *queue) Receive() {}
func (queue *queue) Publish() {}
