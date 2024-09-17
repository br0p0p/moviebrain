package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func main() {
	db = sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
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

type Movie struct {
	BackdropPath  string  `db:"backdrop_path"`
	ID            int     `db:"id"`
	Title         string  `db:"title"`
	OriginalTitle string  `db:"original_title"`
	Overview      string  `db:"overview"`
	PosterPath    string  `db:"poster_path"`
	Popularity    float64 `db:"popularity"`
	ReleaseDate   string  `db:"release_date"`
}

type QueryResults struct {
	Query  string
	Movies []Movie
}

func Search(c echo.Context) error {
	query := c.QueryParam("q")

	if len(query) <= 1 {
		return c.Render(http.StatusOK, "no_results", nil)
	}

	searchTerm := "%" + query + "%"
	// movies := []Movie{}

	fmt.Println("Searching for", searchTerm)

	// err := db.Select(&movies, "SELECT * FROM movie WHERE to_tsvector(title) @@ to_tsquery($1)", query)
	rows, err := db.Queryx("SELECT * FROM movie WHERE title LIKE $1", searchTerm)
	defer rows.Close()

	movies := []map[string]interface{}{}

	// for rows.Next() {
	// 	err = rows.MapScan(movies)
	// 	fmt.Printf("results: %v\n\n", movies)
	// }

	for rows.Next() {
		// Create a map to hold the row data
		rowData := map[string]interface{}{}

		// Use MapScan to scan the row into the map
		if err := rows.MapScan(rowData); err != nil {
			log.Fatal(err)
		}

		// Append the map to the results slice
		movies = append(movies, rowData)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("movies: %v\n\n", movies)

	// Example movie data
	// _movies2 := []map[string]interface{}{
	// 	{"Title": "The Shawshank Redemption", "Director": "Frank Darabont", "Year": 1994},
	// 	{"Title": "The Godfather", "Director": "Francis Ford Coppola", "Year": 1972},
	// 	{"Title": "The Dark Knight", "Director": "Christopher Nolan", "Year": 2008},
	// }

	// Pass the movies list to the template
	data := map[string]interface{}{
		"Query":  query,
		"Movies": movies,
	}

	// b, err := json.Marshal(data)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("JSON:", string(b))

	return c.Render(http.StatusOK, "movies", data)
}
