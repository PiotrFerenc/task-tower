package main

import (
	"github.com/PiotrFerenc/mash2/internal/controllers"
	"github.com/PiotrFerenc/mash2/internal/queues"
	"log"
)

var (
	messageQueue = queues.CreateRabbitMqMessageQueue()
	controller   = controllers.CreateRestController()
)

func main() {
	err := controller.Run("0.0.0.0", "5000")
	if err == nil {
		log.Panic(err)
	}
}
