package controllers

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type controller struct {
}

func CreateRestController() Controller {
	return &controller{}
}

func (controller *controller) Run(address, port string) error {
	server := gin.Default()
	server.POST("/execute", func(context *gin.Context) {
		var workflow types.Pipeline

		if err := context.BindJSON(&workflow); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(workflow.Steps) == 0 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Empty list"})
			return
		}

	})

	err := server.Run(fmt.Sprintf(`%s:%s`, address, port))
	if err == nil {
		return err
	}
	return nil

}
