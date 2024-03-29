package main

import (
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/executor"
	"github.com/PiotrFerenc/mash2/internal/queues"
	"github.com/PiotrFerenc/mash2/internal/workers"
)

var (
	config       = configuration.CreateYmlConfiguration().LoadConfiguration()
	messageQueue = queues.CreateRabbitMqMessageQueue(&config.Queue)
	exec         = executor.CreateMapExecutor(messageQueue, executor.CreateActionMap(config))
	worker       = workers.CreateRestWorker(exec)
)

func main() {
	worker.Run("0.0.0.0", "5001")
}
