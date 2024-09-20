# moviebrain

This is a small app which recommends movies similar to those provided by the user.

This was primarily built to gain more experience with Go and try out HTMX.

https://moviebrain-production.up.railway.app/

## Install

You'll need the following to run the app:

- Go 1.23.1 (or later)
- A running PostgreSQL database
- [`sql-migrate`](https://github.com/rubenv/sql-migrate): `go install github.com/rubenv/sql-migrate/...@latest`
- TMDB API key and bearer token (https://developer.themoviedb.org/v4/docs/getting-started)
- (Optional) `nodemon` is used in development to tighten the iteration loop: `npm i -g nodemon`

## Setup

1. Place your database connection string and TMDB credentials in the `.env` file (see `.env.example` for a template)

2. Run the DB migrations to create the necessary tables:

```sh
sql-migrate up
```

3. Seed the DB with some movie and genre data:

```sh
# Seed all genres
go run ./cmd/dbseed moviegenres
# Seed movies
go run ./cmd/dbseed imdblist
```

4. Run the app:

```sh
PORT=1323 go run .

# or this to re-compile when files change
./dev.sh
```
