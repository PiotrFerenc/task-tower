package main

import (
	"github.com/PiotrFerenc/mash2/src/internal/configuration"
	"github.com/PiotrFerenc/mash2/src/internal/executor"
	"github.com/PiotrFerenc/mash2/src/internal/queues"
	"github.com/PiotrFerenc/mash2/src/internal/workers"
)

var (
	config       = configuration.CreateYmlConfiguration().LoadConfiguration()
	messageQueue = queues.CreateRabbitMqMessageQueue(&config.Queue)
	exec         = executor.CreateMapExecutor(messageQueue)
	worker       = workers.CreateRestWorker(exec)
)

func main() {
	worker.Run("0.0.0.0", "5001")
}
