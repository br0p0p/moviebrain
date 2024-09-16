package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func main() {
	var rootCmd = &cobra.Command{Use: "dbmigrate"}

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the db migrations",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running DB migrations")
			db = sqlx.MustConnect("sqlite3", "data.db")

			schema := `CREATE TABLE genre (
				id INTEGER PRIMARY KEY,
				name VARCHAR(100)
			);`

			result := db.MustExec(schema)
			rowsAffected, err := result.RowsAffected()

			if err != nil {
				panic(err)
			}

			fmt.Println("Rows affected:", rowsAffected)
			fmt.Println("Success.")
		},
	}

	rootCmd.AddCommand(runCmd)
	rootCmd.Execute()

}
