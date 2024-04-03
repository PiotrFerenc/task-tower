package main

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/executor"
	"github.com/PiotrFerenc/mash2/web/persistence"
	"github.com/PiotrFerenc/mash2/web/repositories"
	"github.com/PiotrFerenc/mash2/web/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
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
	e.POST("/pipeline", func(c echo.Context) error {
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

		return c.Render(http.StatusOK, "pipeline-list.html", data)

	})
	e.GET("/", func(c echo.Context) error {
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
	})

	e.GET("/parameters/:action/:id", func(c echo.Context) error {
		idParam := c.Param("id")
		actionName := c.Param("action")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		values := parametersRepository.GetParameters(id)
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		params, ok := parameters[actionName]
		if !ok {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": fmt.Sprintf("Action %s not found", actionName),
			})
		}

		form := mapPropertiesToInputs(params.Inputs(), values)

		data := map[string]interface{}{
			"form": form,
		}
		return c.Render(http.StatusOK, "action-form.html", data)
	})

	e.GET("/pipeline/:id", func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		pipeline, err := pipelineRepository.GetById(id)
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		steps, err := stepsRepository.GetSteps(pipeline.ID)
		if err != nil {
			return c.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
				"error": err.Error(),
			})
		}
		data := map[string]interface{}{
			"Title":    "Strona główna",
			"actions":  parameters,
			"pipeline": pipeline,
			"steps":    steps,
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
func mapPropertiesToInputs(properties []actions.Property, values []types.Parameters) []Input {
	var inputs []Input
	for _, property := range properties {
		value := getParameterValue(values, property.Name)
		input := Input{
			Name:        property.Name,
			Type:        property.Type,
			Description: property.Description,
			Validation:  property.Validation,
			Value:       value.Value,
			Id:          value.ID,
		}
		inputs = append(inputs, input)
	}
	return inputs
}

func getParameterValue(values []types.Parameters, name string) types.Parameters {
	for _, value := range values {
		if value.Key == name {
			return value
		}
	}
	return types.Parameters{}
}

type Input struct {
	Id          uuid.UUID
	Name        string
	Type        string
	Description string
	Validation  string
	Value       string
}
