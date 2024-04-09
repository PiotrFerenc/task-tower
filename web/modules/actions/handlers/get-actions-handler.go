package handlers

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetActionsHandler(parameters map[string]actions.Action) func(c echo.Context) error {
	return func(c echo.Context) error {

		actionName := c.Param("name")
		if actionName == "" {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": "Empty action name",
			})
		}
		actionList := new([]string)
		for name, action := range parameters {
			categoryName := action.GetCategoryName()
			if actionName == categoryName {
				*actionList = append(*actionList, name)
			}
		}
		data := map[string]interface{}{
			"actions": actionList,
		}
		return c.Render(http.StatusOK, "action-new-form.html", data)
	}
}
