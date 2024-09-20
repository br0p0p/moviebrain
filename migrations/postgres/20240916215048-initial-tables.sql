-- +migrate Up
-- Create a table for movie data
CREATE TABLE
  movie (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    original_title VARCHAR(255) NOT NULL,
    overview TEXT,
    backdrop_path VARCHAR(255),
    poster_path VARCHAR(255),
    popularity REAL,
    release_date DATE
  );

-- TODO: Create a GIN index on the "title" and "overview" columns for better full-text search support/performance
-- https://www.postgresql.org/docs/current/textsearch-tables.html
-- CREATE INDEX movie_idx ON movie USING GIN (to_tsvector ('english', title || ' ' || overview));
-- Create the genres table
CREATE TABLE
  genre (id SERIAL PRIMARY KEY, name VARCHAR(100));

-- Create a table to store relationships between movies and genres
CREATE TABLE
  movie_genre (
    movie_id INTEGER REFERENCES movie (id) ON DELETE CASCADE,
    genre_id INTEGER REFERENCES genre (id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, genre_id)
  );

-- +migrate Down
DROP TABLE movie_genre;

DROP TABLE movie;

DROP TABLE genre;