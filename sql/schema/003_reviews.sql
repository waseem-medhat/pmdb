-- +goose Up
CREATE TABLE reviews (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    movie_tmdb_id TEXT NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 0 AND rating <= 10),
    review TEXT NOT NULL,
    public_review INTEGER NOT NULL CHECK (public_review IN (0, 1)),
    UNIQUE(user_id, movie_tmdb_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE reviews;
