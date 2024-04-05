package main

import (
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/executor"
	"github.com/PiotrFerenc/mash2/web/handlers"
	"github.com/PiotrFerenc/mash2/web/persistence"
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

var (
	config               = configuration.CreateYmlConfiguration().LoadConfiguration()
	database             = persistence.CreatePostgresDatabase(&config.Database)
	connection           = database.Connect()
	pipelineRepository   = repositories.CreatePipelineRepository(connection)
	stepsRepository      = repositories.CreateStepsRepository(connection)
	parametersRepository = repositories.CreateParametersRepository(connection)
)

func main() {

	database.RunMigration()

	t := &Template{
		templates: template.Must(template.ParseGlob("web/public/views/*.html")),
	}
	parameters := executor.CreateActionMap(&configuration.Config{})
	e := echo.New()
	e.Renderer = t
	e.Static("/assets", "web/public/static")
	e.GET("/", handlers.CreateHomeHandler(pipelineRepository))
	e.POST("/pipeline", handlers.CreatePipelinesHandler(pipelineRepository))
	e.GET("/pipeline/:id", handlers.CreatePipelineHandler(pipelineRepository, stepsRepository, parameters))
	e.GET("/parameters/:action/:id", handlers.CreateParametersHandler(parametersRepository, parameters))
	e.POST("/parameters", handlers.CreateUpdateParameter(parametersRepository, parameters))
	e.Logger.Fatal(e.Start(":4999"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
