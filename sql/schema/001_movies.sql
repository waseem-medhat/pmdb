-- +goose Up
CREATE TABLE movies (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    year INTEGER,
    created_at TIME NOT NULL,
    updated_at TIME NOT NULL
);

-- +goose Down
DROP TABLE movies;
