package workers

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/internal/executor"
	"github.com/gin-gonic/gin"
)

type worker struct {
	executor executor.Executor
}

func CreateRestWorker(executor executor.Executor) Worker {
	return &worker{
		executor: executor,
	}
}

func (worker *worker) Run(address, port string) {
	server := gin.Default()
	err := server.Run(fmt.Sprintf(`%s:%s`, address, port))
	if err != nil {
		return
	}
}
