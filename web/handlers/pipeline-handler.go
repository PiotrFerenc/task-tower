package handlers

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreatePipelineHandler(pipelineRepository repositories.PipelineRepository, stepsRepository repositories.StepsRepository, parameters map[string]actions.Action) func(c echo.Context) error {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		pipeline, err := pipelineRepository.GetById(id)
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		steps, err := stepsRepository.GetSteps(pipeline.ID)
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		actions := make(map[string][]Action)
		for name, action := range parameters {
			input := Action{
				Outputs:  action.Outputs(),
				Inputs:   action.Inputs(),
				Category: action.GetCategoryName(),
				Name:     name,
			}
			actions[input.Category] = append(actions[input.Category], input)
		}

		data := map[string]interface{}{
			"Title":    "Strona główna",
			"actions":  actions,
			"pipeline": pipeline,
			"steps":    steps,
		}
		return c.Render(http.StatusOK, "index.html", data)
	}
}
func categoryExist(c []string, name string) bool {
	for _, v := range c {
		if v == name {
			return true
		}
	}
	return false
}

type Action struct {
	Outputs  []actions.Property
	Inputs   []actions.Property
	Category string
	Name     string
}
