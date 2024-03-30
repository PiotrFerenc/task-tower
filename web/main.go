package main

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/executor"
	"github.com/PiotrFerenc/mash2/web/persistence"
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
)

var (
	config             = configuration.CreateYmlConfiguration().LoadConfiguration()
	database           = persistence.CreatePostgresDatabase(&config.Database)
	connection         = database.Connect()
	pipelineRepository = repositories.CreatePipelineRepository(connection)
)

func main() {

	database.RunMigration()

	t := &Template{
		templates: template.Must(template.ParseGlob("web/public/views/*.html")),
	}
	e := echo.New()
	e.Renderer = t
	e.Static("/assets", "web/public/static")
	e.GET("/", func(c echo.Context) error {
		pipelines, err := pipelineRepository.GetAll()
		if err != nil {
			return c.Render(http.StatusOK, "pipelines.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		data := map[string]interface{}{
			"Title":     "Strona główna",
			"pipelines": pipelines,
		}
		return c.Render(http.StatusOK, "pipelines.html", data)
	})
	e.GET("/pipeline/:id", func(c echo.Context) error {
		data := map[string]interface{}{
			"Title":   "Strona główna",
			"actions": mapItems(executor.CreateActionMap(&configuration.Config{})),
		}
		return c.Render(http.StatusOK, "index.html", data)
	})
	e.Logger.Fatal(e.Start(":4999"))

}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func mapItems(actions map[string]actions.Action) map[string]ActionParameter {
	result := make(map[string]ActionParameter)
	for key, action := range actions {
		result[key] = ActionParameter{
			Inputs:  action.Inputs(),
			Outputs: action.Outputs(),
		}
	}
	return result

}

type ActionParameter struct {
	Inputs  []actions.Property
	Outputs []actions.Property
}
