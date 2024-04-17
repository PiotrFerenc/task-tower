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

// Run is a method that starts the server and listens for incoming requests.
// It takes an address and a port as parameters and returns an error if any
// errors occur during the server startup.
//
// The method first creates a new Gin server instance with default middleware
// using gin.Default(). It then registers two routes: one for GET requests to
// "/process/:id" and one for POST requests to "/execute". Both routes use
// handler functions to handle the requests.
//
// Next, it attempts to run the server by calling server.Run with the provided
// address and port. If an error occurs during the server startup, it returns
// the error. Otherwise, it returns nil.
func (controller *controller) Run(address, port string) error {
	server := gin.Default()
	server.GET("/process/:id", getByIdHandler(controller))
	server.POST("/execute", executeHandler(controller))

	err := server.Run(fmt.Sprintf(`%s:%s`, address, port))
	if err == nil {
		return err
	}
	return nil

}

// getByIdHandler is a handler function that retrieves a process by its ID and returns
// the status of the process.
// It takes a controller instance as a parameter and returns a function that has a
// gin.Context as its parameter.
// The function first retrieves the ID parameter from the gin.Context, parses it into a UUID,
// and checks if the ID format is valid.
// If the ID format is invalid, it returns a JSON response with a 400 Bad Request status
// and an error message.
// If the ID format is valid, it calls the GetById method of the processRepository
// to retrieve the process based on the ID.
// If the process is not found, it returns a JSON response with a 404 Not Found status
// and an error message.
// If the process is found, it returns a JSON response with a 200 OK status
// and the process status.
func getByIdHandler(ctrl *controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		parameterId := c.Param("id")
		id, err := uuid.Parse(parameterId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		process, err := ctrl.processRepository.GetById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Process not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": process.Status})
	}
}

// executeHandler is a handler function that takes a controller instance as a parameter and returns a function
// that has a gin.Context as its parameter. The function first binds the JSON payload from the context to a
// types.Pipeline struct. If there is an error in binding, the function aborts with a 400 Bad Request status
// and returns the error message. Next, it validates the pipeline using the Validate method of the
// types.Pipeline struct. If there is an error in validation, the function aborts with a 400 Bad Request status
// and returns the error message. It checks if the pipeline contains any tasks, and if not, aborts with a
// 400 Bad Request status and returns an error message from the Message.EmptyTaskList constant. If the pipeline
// is valid, it calls the Run method of the pipelineService to execute the pipeline and gets the unique process ID.
// If there is an error in executing the pipeline, it aborts with a 400 Bad Request status and returns the error.
// Finally, it returns a JSON response with a 200 OK status and the process ID.
func executeHandler(c *controller) func(context *gin.Context) {
	return func(context *gin.Context) {

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

		processId, err := c.pipelineService.Run(&pipe)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		}
		context.JSON(http.StatusOK, gin.H{"processId": processId})
	}
}
