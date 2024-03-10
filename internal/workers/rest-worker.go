package workers

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

type worker struct {
}

func CreateRestWorker() Worker {
	return &worker{}
}

func (worker *worker) Run(address, port string) {
	server := gin.Default()
	server.POST("/execute/:action", func(context *gin.Context) {
		name := context.Param("action")
		var parameters actions.ActionContext

		if err := context.BindJSON(&parameters); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := internal.Executor(name, parameters)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		}
	})

	err := server.Run(fmt.Sprintf(`%s:%s`, address, port))
	if err != nil {
		return
	}
}
