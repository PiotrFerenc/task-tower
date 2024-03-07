package main

import (
	workers2 "github.com/PiotrFerenc/mash2/internal/workers"
)

func main() {
	config := &workers2.RestWorker{
		Address: "0.0.0.0",
		Port:    "5000",
	}

	worker := CreateWorker(config)
	worker.Run()
}

func CreateWorker(worker workers2.Worker) workers2.Worker {
	return worker
}
