package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.POST("/execute", func(context *gin.Context) {
		var workflow Workflow

		if err := context.BindJSON(&workflow); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(workflow.Steps) == 0 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Empty list"})
			return
		}

	})

	err := server.Run(fmt.Sprintf(`%s:%s`, "0.0.0.0", "5000"))
	if err != nil {
		return
	}
}

type Workflow struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"`
}
