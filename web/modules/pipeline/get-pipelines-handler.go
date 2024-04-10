package pipeline

import (
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetPipelinesHandler(pipelineRepository repositories.PipelineRepository) func(c echo.Context) error {
	return func(c echo.Context) error {
		pipelines, err := pipelineRepository.GetAll()
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		data := map[string]interface{}{
			"Title":     "Strona główna",
			"pipelines": pipelines,
		}
		return c.Render(http.StatusOK, "pipelines.html", data)
	}
}
