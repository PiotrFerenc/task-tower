package modules

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/controllers"
	"github.com/PiotrFerenc/mash2/web/modules/actions/handlers"
	"github.com/PiotrFerenc/mash2/web/modules/editor"
	"github.com/PiotrFerenc/mash2/web/modules/parameters"
	"github.com/PiotrFerenc/mash2/web/modules/pipeline"
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/labstack/echo/v4"
)

func RegisterActionModule(e *echo.Echo, actions map[string]actions.Action, stepsRepository repositories.StepsRepository, parametersRepository repositories.ParametersRepository) (g *echo.Group) {
	api := e.Group("/actions")
	api.GET("/categories", handlers.GetCategoriesHandler(actions))
	api.GET("/:name", handlers.GetActionsHandler(actions))
	api.GET("/:name/parameters", handlers.GetActionParametersHandler(actions))
	api.POST("/:name/:pipelineId", handlers.AddActionToPipelineHandler(stepsRepository, parametersRepository))
	return api
}

func RegisterDashboardModule(e *echo.Echo, client controllers.ControllerClient) {
	e.GET("/", editor.CreateEditorHandler())
	e.POST("/execute", editor.ExecutePipelineHandler(client))
}

func RegisterParametersModule(e *echo.Echo, parametersRepository repositories.ParametersRepository, actions map[string]actions.Action) (g *echo.Group) {
	api := e.Group("/parameters")
	e.GET("/parameters/:action/:id", parameters.GetParametersHandler(parametersRepository, actions))
	e.POST("/parameters", parameters.UpdateParameter(parametersRepository, actions))
	return api
}

func RegisterPipelineModule(e *echo.Echo, pipelineRepository repositories.PipelineRepository, stepsRepository repositories.StepsRepository, parameters map[string]actions.Action) (g *echo.Group) {
	api := e.Group("/pipeline")
	e.GET("/pipeline", pipeline.GetPipelinesHandler(pipelineRepository))
	api.POST("/", pipeline.CreatePipelinesHandler(pipelineRepository))
	api.GET("/:id", pipeline.GetPipelineHandler(pipelineRepository, stepsRepository, parameters))
	return api
}
func RegisterEditor(e *echo.Echo, pipelineRepository repositories.PipelineRepository, stepsRepository repositories.StepsRepository, parameters map[string]actions.Action, parametersRepository repositories.ParametersRepository) {
	RegisterPipelineModule(e, pipelineRepository, stepsRepository, parameters)
	RegisterParametersModule(e, parametersRepository, parameters)
	RegisterActionModule(e, parameters, stepsRepository, parametersRepository)
}
