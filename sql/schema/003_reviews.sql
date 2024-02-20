-- +goose Up
CREATE TABLE reviews (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  movie_id TEXT NOT NULL,
  rating INTEGER CHECK (rating >= 1 AND rating <= 5),
  review TEXT,
  created_at TIME NOT NULL,
  updated_at TIME NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (movie_id) REFERENCES movies(id)
);


-- +goose Down
DROP TABLE reviews;
