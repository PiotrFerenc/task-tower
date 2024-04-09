package main

import (
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/web/modules"
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
	t := LoadTemplates([]string{
		"web/public/views/*.html",
		"web/public/views/dashboard/*.html",
	})
	//	parameters := executor.CreateActionMap(&configuration.Config{})
	e := echo.New()
	e.Renderer = t

	e.Static("/assets", "web/public/static")

	modules.RegisterDashboardModule(e, pipelineRepository)
	modules.RegisterEditor(e, pipelineRepository)
	e.Logger.Fatal(e.Start(":4999"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func LoadTemplates(patterns []string) *Template {
	var tmpl *template.Template

	for _, pattern := range patterns {
		if tmpl == nil {
			tmpl = template.Must(template.ParseGlob(pattern))
		} else {
			tmpl = template.Must(tmpl.ParseGlob(pattern))
		}
	}
	return &Template{templates: tmpl}
}
