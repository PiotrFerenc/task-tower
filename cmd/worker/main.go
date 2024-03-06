package main

import "github.com/PiotrFerenc/mash2/cmd/worker/workers"

func main() {
	config := &workers.RestWorker{
		Address: "0.0.0.0",
		Port:    "5000",
	}

	worker := CreateWorker(config)
	worker.Run()
}

func CreateWorker(worker workers.Worker) workers.Worker {
	return worker
}
