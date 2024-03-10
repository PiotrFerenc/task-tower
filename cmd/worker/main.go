package main

import (
	"github.com/PiotrFerenc/mash2/internal/executor"
	"github.com/PiotrFerenc/mash2/internal/workers"
)

var (
	exec   = executor.CreateMapExecutor()
	worker = workers.CreateRestWorker(exec)
)

func main() {
	worker.Run("0.0.0.0", "5001")
}
