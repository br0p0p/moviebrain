# moviebrain

This is a small app which recommends movies similar to those provided by the user.

This was primarily built to gain more experience with Go and try out HTMX.

https://moviebrain-production.up.railway.app/

## Install

You'll need the following to run the app:

- Go 1.23.1 (or later)
- A running PostgreSQL database
- TMDB API key and token (https://developer.themoviedb.org/v4/docs/getting-started)
- (Optional) `nodemon` is used in development to tighten the iteration loop: `npm i -g nodemon`

## Setup

- Place your database connection string and TMDB credentials in the `.env` file (see `.env.example` for a template)
- Run the DB migrations to create the necessary tables: `sql-migrate up`
- Seed the DB with some movie and genre data: `go run ./cmd/dbseed moviegenres` and `go run ./cmd/dbseed imdblist`
- Run the app: `PORT=1323 go run .`
