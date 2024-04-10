package handlers

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/web/modules/pipeline"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetActionParametersHandler(parameters map[string]actions.Action) func(c echo.Context) error {
	return func(c echo.Context) error {
		actionName := c.Param("name")
		if actionName == "" {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": "Empty action name",
			})
		}
		actionList := new([]pipeline.Action)
		for name, action := range parameters {
			if actionName == name {
				input := pipeline.Action{
					Outputs:  action.Outputs(),
					Inputs:   action.Inputs(),
					Category: action.GetCategoryName(),
					Name:     name,
				}
				*actionList = append(*actionList, input)
			}
		}
		data := map[string]interface{}{
			"actions": actionList,
		}
		return c.Render(http.StatusOK, "action-parameters.html", data)
	}
}
