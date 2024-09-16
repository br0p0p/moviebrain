package main

import (
	"fmt"
	"moviebrain/moviebrain"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var db *sqlx.DB

func main() {

	var rootCmd = &cobra.Command{Use: "dbseed"}

	var movieGenresCmd = &cobra.Command{
		Use:   "moviegenres",
		Short: "Seed the database with all available movie genres",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running movie genres seed")
			db = sqlx.MustConnect("sqlite3", "data.db")
			tmdbClient := moviebrain.GetTmdbV3Client()
			result, err := tmdbClient.GetGenreMovieList(nil)
			if err != nil {
				panic(err)
			}

			fmt.Println("First item:", result.Genres[0])

			var totalRowsAffected int64 = 0
			genreInsert := `INSERT INTO genre (id, name) VALUES (?, ?)`

			for _, v := range result.Genres {
				result := db.MustExec(genreInsert, v.ID, v.Name)
				rowsAffected, _ := result.RowsAffected()
				totalRowsAffected += rowsAffected
			}

			fmt.Println("Rows affected:", totalRowsAffected)
			fmt.Println("Success.")
		},
	}

	var imdblistCmd = &cobra.Command{
		Use:   "imdblist",
		Short: "Seed the database with the top 250 IMDB movies",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running imdblist seed")
			db = sqlx.MustConnect("sqlite3", "data.db")
			tmdbClient := moviebrain.GetTmdbV4Client()
			result, err := tmdbClient.GetListDetails(moviebrain.IMDB_TOP_250_ID, nil)
			if err != nil {
				panic(err)
			}

			fmt.Println("First item:", result.Items[0])

			var totalRowsAffected int64 = 0

			fmt.Println("Rows affected:", totalRowsAffected)
			fmt.Println("Success.")
		},
	}

	rootCmd.AddCommand(movieGenresCmd)
	rootCmd.AddCommand(imdblistCmd)
	rootCmd.Execute()

}
