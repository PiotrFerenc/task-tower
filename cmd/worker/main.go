package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func main() {
	config := &RestWorker{
		Address: "0.0.0.0",
		Port:    "5000",
	}

	worker := CreateWorker(config)
	worker.Run()
}

func CreateWorker(worker Worker) Worker {
	return worker
}

type Worker interface {
	Run()
}

type RestWorker struct {
	Address string
	Port    string
}

func (worker RestWorker) Run() {
	server := gin.Default()
	server.POST("/execute/:action", func(context *gin.Context) {
		name := context.Param("action")
		var parameters ActionContext

		if err := context.BindJSON(&parameters); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := Executor(name, parameters)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		}
	})

	err := server.Run(fmt.Sprintf(`%s:%s`, worker.Address, worker.Port))
	if err != nil {
		return
	}
}

func Executor(actionName string, parameters ActionContext) error {
	actions := map[string]Action{
		"hallo": Hallo{},
		"sleep": Sleep{},
	}
	action, exist := actions[actionName]
	if !exist {
		return fmt.Errorf("action %v not found", actionName)
	}
	action.Execute(parameters)
	return nil
}

type Action interface {
	Execute(parameters ActionContext) string
}

type Hallo struct {
}

func (receiver Hallo) Execute(parameters ActionContext) string {
	name := parameters.Parameters["name"]

	msg := "Hallo " + name
	log.Print(msg)
	return msg
}

type Sleep struct {
}

func (receiver Sleep) Execute(parameters ActionContext) string {

	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		log.Print("sleep")
	}
	return ""
}

type ActionContext struct {
	Parameters map[string]string `json:"parameters"`
}
