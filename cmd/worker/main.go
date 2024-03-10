package main

import (
	"github.com/PiotrFerenc/mash2/internal/workers"
)

var (
	worker = workers.CreateRestWorker()
)

func main() {
	worker.Run("0.0.0.0", "5001")
}
