package editor

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateEditorHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		data := map[string]interface{}{
			"Title": "Edytor",
		}
		return c.Render(http.StatusOK, "edytor.html", data)
	}
}
