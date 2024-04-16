package controllers

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/api/types"
	Message "github.com/PiotrFerenc/mash2/internal/consts"
	"github.com/PiotrFerenc/mash2/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type controller struct {
	pipelineService services.PipelineService
}

func CreateRestController(pipelineService services.PipelineService) Controller {
	return &controller{
		pipelineService: pipelineService,
	}
}

func (controller *controller) Run(address, port string) error {
	server := gin.Default()
	server.POST("/execute", func(context *gin.Context) {
		var pipe types.Pipeline

		if err := context.BindJSON(&pipe); err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := pipe.Validate(); err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(pipe.Tasks) == 0 {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": Message.EmptyTaskList})
			return
		}

		processId, err := controller.pipelineService.Run(&pipe)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		}
		context.JSON(http.StatusOK, gin.H{"processId": processId})
	})

	err := server.Run(fmt.Sprintf(`%s:%s`, address, port))
	if err == nil {
		return err
	}
	return nil

}
