package main

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
)

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e := echo.New()
	e.Renderer = t
	e.Static("/css", "public/static")
	e.GET("/", func(c echo.Context) error {
		data := map[string]interface{}{
			"Title":   "Strona główna",
			"Header":  "Witaj w Echo",
			"Message": "To jest strona główna wygenerowana z szablonu.",
		}
		return c.Render(http.StatusOK, "index.html", data)
	})
	e.Logger.Fatal(e.Start(":80"))

}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
