package pipeline

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetPipelineHandler(pipelineRepository repositories.PipelineRepository, stepsRepository repositories.StepsRepository, parameters map[string]actions.Action) func(c echo.Context) error {
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
		actionList := make(map[string][]Action)
		for name, action := range parameters {
			input := Action{
				Outputs:  action.Outputs(),
				Inputs:   action.Inputs(),
				Category: action.GetCategoryName(),
				Name:     name,
			}
			actionList[input.Category] = append(actionList[input.Category], input)
		}

		data := map[string]interface{}{
			"Title":      "Strona główna",
			"categories": actionList,
			"pipeline":   pipeline,
			"steps":      steps,
		}
		return c.Render(http.StatusOK, "index.html", data)
	}
}

type Action struct {
	Outputs  []actions.Property
	Inputs   []actions.Property
	Category string
	Name     string
}
