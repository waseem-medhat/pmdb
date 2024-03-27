-- +goose Up
ALTER TABLE users
    ADD COLUMN bio TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE users
    DROP COLUMN bio;

