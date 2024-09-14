package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Static("/assets", "assets")

	tmpl := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e.Renderer = tmpl

	e.GET("/", Index)
	e.GET("/hello", Hello)
	e.GET("/search", Search)

	addr := ":" + getPort()
	e.Logger.Fatal(e.Start(addr))
}

func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Please specify the HTTP port as environment variable, e.g. env PORT=8081 go run http-server.go")
	}

	return port
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
	return c.Render(http.StatusOK, "base", "World")
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

type User struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

func Search(c echo.Context) error {
	query := c.QueryParam("q")

	// Example movie data
	movies := []map[string]interface{}{
		{"Title": "The Shawshank Redemption", "Director": "Frank Darabont", "Year": 1994},
		{"Title": "The Godfather", "Director": "Francis Ford Coppola", "Year": 1972},
		{"Title": "The Dark Knight", "Director": "Christopher Nolan", "Year": 2008},
	}

	// Pass the movies list to the template
	data := map[string]interface{}{
		"Query":  query,
		"Movies": movies,
	}

	return c.Render(http.StatusOK, "movies", data)
}
