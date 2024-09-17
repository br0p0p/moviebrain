-- +migrate Up
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

CREATE TABLE
  genre (id SERIAL PRIMARY KEY, name VARCHAR(100));

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