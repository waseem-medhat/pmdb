-- +goose Up
CREATE TABLE reviews (
  id TEXT PRIMARY KEY,
  user_id INTEGER NOT NULL,
  movie_tmdb_id TEXT NOT NULL,
  rating INTEGER CHECK (rating >= 0 AND rating <= 10),
  review TEXT,
  public_review INTEGER NOT NULL CHECK (public_review IN (0, 1)),
  FOREIGN KEY (user_id) REFERENCES users(id)
);


-- +goose Down
DROP TABLE reviews;
