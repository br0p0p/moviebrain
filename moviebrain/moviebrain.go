package moviebrain

import (
	"fmt"
	"os"

	tmdb "github.com/cyruzin/golang-tmdb"
)

const IMDB_TOP_250_ID = 634

func GetTmdbV3Client() *tmdb.Client {
	tmdbClient, err := tmdb.Init(os.Getenv("TMDB_APIKEY"))
	if err != nil {
		fmt.Println(err)
	}

	tmdbClient.SetClientAutoRetry()

	return tmdbClient
}

func GetTmdbV4Client() *tmdb.Client {
	// Using v4
	tmdbClient, err := tmdb.InitV4(os.Getenv("TMDB_BEARER_TOKEN"))
	if err != nil {
		fmt.Println(err)
	}

	tmdbClient.SetClientAutoRetry()

	return tmdbClient
}

func getSeedMovies() {
	tmdbClient := GetTmdbV4Client()

	listDetails, err := tmdbClient.GetListDetails(IMDB_TOP_250_ID, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(listDetails.Items[0])
}
