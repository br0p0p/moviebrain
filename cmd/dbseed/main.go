package main

import (
	"fmt"
	"moviebrain/moviebrain"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

func main() {
	tmdbClient := moviebrain.GetTmdbV4Client()

	var rootCmd = &cobra.Command{Use: "dbseed"}

	var imdblistCmd = &cobra.Command{
		Use:   "imdblist",
		Short: "Seed the database with the top 250 IMDB movies",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running imdblist seed")
			listDetails, err := tmdbClient.GetListDetails(moviebrain.IMDB_TOP_250_ID, nil)
			if err != nil {
				panic(err)
			}

			fmt.Println("First item:", listDetails.Items[0])
		},
	}

	rootCmd.AddCommand(imdblistCmd)
	rootCmd.Execute()

}
