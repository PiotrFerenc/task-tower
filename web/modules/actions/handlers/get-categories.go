package handlers

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetCategoriesHandler(parameters map[string]actions.Action) func(c echo.Context) error {
	return func(c echo.Context) error {
		categories := make([]string, len(parameters))
		for _, action := range parameters {
			categoryName := action.GetCategoryName()
			categories = append(categories, categoryName)
		}
		data := map[string]interface{}{
			"categories": categories,
		}
		return c.Render(http.StatusOK, "actions-category.html", data)
	}
}
