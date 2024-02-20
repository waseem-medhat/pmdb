-- +goose Up
CREATE TABLE movies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    year INTEGER
);

-- +goose Down
DROP TABLE movies;
