package controllers

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/api/types"
	Message "github.com/PiotrFerenc/mash2/internal/consts"
	"github.com/PiotrFerenc/mash2/internal/repositories"
	"github.com/PiotrFerenc/mash2/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type controller struct {
	pipelineService   services.PipelineService
	processRepository repositories.ProcessRepository
}

func CreateRestController(pipelineService services.PipelineService, repository repositories.ProcessRepository) Controller {
	return &controller{
		pipelineService:   pipelineService,
		processRepository: repository,
	}
}

func (controller *controller) Run(address, port string) error {
	server := gin.Default()
	server.GET("/process/:id", func(c *gin.Context) {
		parameterId := c.Param("id")
		id, err := uuid.Parse(parameterId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		process, err := controller.processRepository.GetById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Process not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": process.Status})
	})
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
