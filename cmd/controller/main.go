package main

import (
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/controllers"
	"github.com/PiotrFerenc/mash2/internal/queues"
	"github.com/PiotrFerenc/mash2/internal/repositories"
	"github.com/PiotrFerenc/mash2/internal/services"
	"log"
)

var (
	config            = configuration.CreateYmlConfiguration().LoadConfiguration()
	messageQueue      = queues.CreateRabbitMqMessageQueue(&config.Queue)
	processRepository = repositories.CreateMongoDbRepository(&config.Database)
	processService    = services.CreateProcessService(processRepository)
	pipeLineService   = services.CreatePipelineService(messageQueue, processService, services.CreateOnSuccessFunc(), services.CreateOnFailFunc(), services.CreateOnFinishFunc())
	controller        = controllers.CreateRestController(pipeLineService, processRepository)
)

func main() {
	err := controller.Run("0.0.0.0", "5000")
	if err == nil {
		log.Panic(err)
	}
}
