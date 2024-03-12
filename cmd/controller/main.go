package main

import (
	"github.com/PiotrFerenc/mash2/internal/controllers"
	"github.com/PiotrFerenc/mash2/internal/queues"
	"github.com/PiotrFerenc/mash2/internal/repositories"
	"github.com/PiotrFerenc/mash2/internal/services"
	"log"
)

var (
	messageQueue      = queues.CreateRabbitMqMessageQueue()
	processRepository = repositories.CreatePostgresRepository()
	processService    = services.CreateProcessService(processRepository)
	pipeLineService   = services.CreatePipelineService(messageQueue, processService)
	controller        = controllers.CreateRestController(pipeLineService)
)

func main() {
	err := controller.Run("0.0.0.0", "5000")
	if err == nil {
		log.Panic(err)
	}
}
