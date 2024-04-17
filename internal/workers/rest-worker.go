package workers

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/internal/executor"
	"github.com/gin-gonic/gin"
)

type worker struct {
	executor executor.Executor
}

// CreateRestWorker creates a new RestWorker with the given executor.
// It takes an executor.Executor as a parameter and returns a Worker.
func CreateRestWorker(executor executor.Executor) Worker {
	return &worker{
		executor: executor,
	}
}

// Run starts the worker by running a Gin server on the specified address and port.
// It takes the address and port as parameters and does not return any value.
// If an error occurs while starting the server, the method returns without any action.
func (worker *worker) Run(address, port string) {
	server := gin.Default()
	err := server.Run(fmt.Sprintf(`%s:%s`, address, port))
	if err != nil {
		return
	}
}
