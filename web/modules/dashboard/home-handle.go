package dashboard

import (
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateHomeHandler(pipelineRepository repositories.PipelineRepository) func(c echo.Context) error {
	return func(c echo.Context) error {

		data := map[string]interface{}{
			"Title": "Strona główna",
		}
		return c.Render(http.StatusOK, "home.html", data)
	}
}
