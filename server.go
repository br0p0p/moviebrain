package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func main() {
	// The app will halt at this step if no database connection can be made.
	// In a real world app we would probably want to serve an error page in this case instead of completely halting the server.
	db = sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Static("/assets", "assets")

	tmpl := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e.Renderer = tmpl

	e.GET("/", Index)
	e.GET("/search", Search)

	addr := ":" + getPort()
	e.Logger.Fatal(e.Start(addr))
}

func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Please specify the HTTP port as environment variable, e.g. `PORT=8081 go run http-server.go`")
	}

	return port
}

type (
	MovieSuggestionsForm struct {
		ids []string `json:"movieIds" validate:"required"`
	}
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "base", nil)
}

// This type contains a subset of the fields available to us via TMDB. In a real world scenario, this would probably have many more fields from multiple sources.
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

// Since we're using HTMX on the client, it expects AJAX responses to be HTML which it will insert into the DOM.
// If we were using something like React on the frontend, this endpoint would probably return JSON.
func Search(c echo.Context) error {
	query := c.QueryParam("q")

	// To avoid unwanted server load, don't even query the DB if there aren't enough characters in the query.
	if len(query) <= 1 {
		return c.Render(http.StatusOK, "no_results", nil)
	}

	fmt.Println("Searching for", query)

	// TODO: Use PostgreSQL's builtin full-text search to significantly improve search results quality
	// rows, err := db.Queryx("SELECT * FROM movie WHERE to_tsvector(title) @@ to_tsquery($1)", query)
	rows, err := db.Queryx("SELECT * FROM movie WHERE title LIKE $1", "%"+query+"%")
	defer rows.Close()

	// Create a slice which will hold all `Movie` results
	movies := []map[string]interface{}{}

	for rows.Next() {
		// Create a map to hold the row data
		// There's probably a better way to do this so that `rowData` is properly typed as `Movie`
		rowData := map[string]interface{}{}

		// Use MapScan to scan the row into the map
		if err := rows.MapScan(rowData); err != nil {
			log.Fatal(err)
		}

		rowData["full_poster_path"] = tmdb.GetImageURL(rowData["poster_path"].(string), tmdb.W154)

		// Append the movie to the `movies` slice
		movies = append(movies, rowData)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Pass the movies list to the template
	data := map[string]interface{}{
		"Query":  query,
		"Movies": movies,
	}

	return c.Render(http.StatusOK, "movies", data)
}
