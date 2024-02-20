-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIME NOT NULL,
    updated_at TIME NOT NULL
);

-- +goose Down
DROP TABLE users;
