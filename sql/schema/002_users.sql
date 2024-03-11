-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    user_name TEXT NOT NULL,
    display_name text NOT NULL
);

-- +goose Down
DROP TABLE users;
