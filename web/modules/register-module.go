package modules

import (
	"github.com/PiotrFerenc/mash2/web/modules/dashboard"
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/labstack/echo/v4"
)

func RegisterActionModule(e *echo.Echo) (g *echo.Group) {
	api := e.Group("/actions")
	return api
}

func RegisterDashboardModule(e *echo.Echo, pipelineRepository repositories.PipelineRepository) {
	e.GET("/", dashboard.CreateHomeHandler(pipelineRepository))
}

func RegisterParametersModule(e *echo.Echo) (g *echo.Group) {
	api := e.Group("/parameters")
	return api
}

func RegisterPipelineModule(e *echo.Echo, pipelineRepository repositories.PipelineRepository) (g *echo.Group) {
	api := e.Group("/pipeline")
	api.GET("/pipelines", dashboard.CreateHomeHandler(pipelineRepository))
	return api
}
func RegisterEditor(e *echo.Echo, pipelineRepository repositories.PipelineRepository) {
	RegisterPipelineModule(e, pipelineRepository)
	RegisterParametersModule(e)
	RegisterActionModule(e)
}
