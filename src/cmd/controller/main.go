package main

import (
	"github.com/PiotrFerenc/mash2/src/internal/configuration"
	"github.com/PiotrFerenc/mash2/src/internal/controllers"
	"github.com/PiotrFerenc/mash2/src/internal/queues"
	"github.com/PiotrFerenc/mash2/src/internal/repositories"
	"github.com/PiotrFerenc/mash2/src/internal/services"
	"log"
)

var (
	config            = configuration.CreateYmlConfiguration().LoadConfiguration()
	messageQueue      = queues.CreateRabbitMqMessageQueue(&config.Queue)
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
