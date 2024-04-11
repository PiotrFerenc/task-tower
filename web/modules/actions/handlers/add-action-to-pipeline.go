package handlers

import (
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddActionToPipelineHandler(stepRepository repositories.StepsRepository, parameterRepository repositories.ParametersRepository) func(ss echo.Context) error {
	return func(ss echo.Context) error {
		actionName := ss.Param("name")
		if actionName == "" {
			return ss.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": "Empty action name",
			})
		}
		pipelineId := ss.Param("pipelineId")
		id, err := uuid.Parse(pipelineId)
		if err != nil {
			return ss.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		actionsFromRequestForm := make(map[string]interface{})

		err = ss.Bind(&actionsFromRequestForm)
		if err != nil {
			return ss.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": "Empty action name",
			})
		}
		stepId, err := stepRepository.Save(actionName, id)
		if err != nil {
			return ss.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		err = parameterRepository.AddParameters(stepId, actionsFromRequestForm)
		if err != nil {
			return ss.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		return nil
	}
}
