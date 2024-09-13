package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e.Renderer = t

	e.GET("/", Index)

	e.GET("/hello", Hello)

	e.Logger.Fatal(e.Start(":1323"))
}

type (
	MovieSuggestionsForm struct {
		ids []string `json:"movieIds" validate:"required"`
	}
)

// func postMovieSelections(c echo.Context) error {

// }

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}
