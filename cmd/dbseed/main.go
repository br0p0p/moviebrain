// A basic CLI to make database seeding easier

package main

import (
	"fmt"
	"moviebrain/moviebrain"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var db *sqlx.DB

func main() {

	var rootCmd = &cobra.Command{Use: "dbseed"}

	// Queries all available genres from TMDB, inserts them into the `genre` table
	var movieGenresCmd = &cobra.Command{
		Use:   "moviegenres",
		Short: "Seed the database with all available movie genres",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running movie genres seed")
			db = sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))
			tmdbClient := moviebrain.GetTmdbV3Client()
			result, err := tmdbClient.GetGenreMovieList(nil)
			if err != nil {
				panic(err)
			}

			var totalRowsAffected int64 = 0
			genreInsert := `INSERT INTO genre (id, name) VALUES ($1, $2)`

			for _, v := range result.Genres {
				result := db.MustExec(genreInsert, v.ID, v.Name)
				rowsAffected, _ := result.RowsAffected()
				totalRowsAffected += rowsAffected
			}

			fmt.Println("Rows affected:", totalRowsAffected)
			fmt.Println("Success.")
		},
	}

	// Fetches movies from the "IMDB Top 250" list on TMDB, inserts them into the `movie` table, and creates relations for the attached genres.
	var imdblistCmd = &cobra.Command{
		Use:   "imdblist",
		Short: "Seed the database with the IMDB Top 250 movies",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running imdblist seed")
			db = sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))
			tmdbClient := moviebrain.GetTmdbV4Client()

			var totalMovieRowsAffected int64 = 0
			var totalMovieGenreRowsAffected int64 = 0
			movieInsert := `INSERT INTO movie (id, title, original_title, overview, backdrop_path, poster_path, popularity, release_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
			movieGenreInsert := `INSERT INTO movie_genre (movie_id, genre_id) VALUES ($1, $2)`

			// There are only 20 results per page so we have to paginate data queries
			var page int = 1

			for {
				fmt.Println("Fetching page", strconv.Itoa(page))
				result, err := tmdbClient.GetListDetails(moviebrain.IMDB_TOP_250_ID, map[string]string{
					"page": strconv.Itoa(page),
				})

				if err != nil {
					panic(err)
				}

				genreMatches := 0

				// Iterate over all the movie results, inserting them into the `movie` table
				for _, m := range result.Items {
					result := db.MustExec(movieInsert, m.ID, m.Title, m.OriginalTitle, m.Overview, m.BackdropPath, m.PosterPath, m.Popularity, m.ReleaseDate)
					movieRowsAffected, _ := result.RowsAffected()
					totalMovieRowsAffected += movieRowsAffected

					// For every genre
					for _, g := range m.GenreIDs {
						result := db.MustExec(movieGenreInsert, m.ID, g)

						movieGenreRowsAffected, _ := result.RowsAffected()
						totalMovieGenreRowsAffected += movieGenreRowsAffected
					}

					genreMatches += len(m.GenreIDs)

				}

				// This is a quick and dirty way to check if we're on the last page.
				// Since there are max 20 results per page, if the current page has less than 20 items we can (semi)safely assume we're on the last page.
				if len(result.Items) < 20 {
					break
				}

				page += 1
			}

			fmt.Println("Inserted", totalMovieRowsAffected, "into table 'movie'")
			fmt.Println("Inserted", totalMovieGenreRowsAffected, "into table 'movie_genre'")
			fmt.Println("Inserted", totalMovieRowsAffected+totalMovieGenreRowsAffected, "rows total")

			fmt.Println("Success.")
		},
	}

	rootCmd.AddCommand(movieGenresCmd)
	rootCmd.AddCommand(imdblistCmd)
	rootCmd.Execute()

}
