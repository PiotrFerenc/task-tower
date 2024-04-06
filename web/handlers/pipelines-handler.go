package handlers

import (
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreatePipelinesHandler(pipelineRepository repositories.PipelineRepository) func(ctx echo.Context) error {
	return func(c echo.Context) error {
		data := map[string]interface{}{}
		if err := c.Bind(&data); err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		name, ok := data["pipeline-name"]
		if !ok {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": "pipeline-name is required",
			})
		}
		pipelineName := name.(string)
		pipeline, err := pipelineRepository.Save(pipelineName)
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		tmp := map[string]interface{}{
			"ID":   pipeline.ID,
			"Name": pipeline.Name,
		}

		return c.Render(http.StatusOK, "pipeline-list.html", tmp)

	}

}
