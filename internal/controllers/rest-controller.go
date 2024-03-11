package controllers

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type controller struct {
	processService services.ProcessService
}

func CreateRestController(processService services.ProcessService) Controller {
	return &controller{
		processService: processService,
	}
}

func (controller *controller) Run(address, port string) error {
	server := gin.Default()
	server.POST("/execute", func(context *gin.Context) {
		var pipe types.Pipeline

		if err := context.BindJSON(&pipe); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(pipe.Stages) == 0 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Empty list"})
			return
		}

	})

	err := server.Run(fmt.Sprintf(`%s:%s`, address, port))
	if err == nil {
		return err
	}
	return nil

}
