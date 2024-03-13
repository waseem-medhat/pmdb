-- name: GetUser :one
SELECT id, display_name, user_name FROM users
WHERE user_name = ? LIMIT 1;

-- name: GetUserForLogin :one
SELECT user_name, password FROM users
WHERE user_name = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY display_name;

-- name: CreateUser :one
INSERT INTO users ( id, user_name, display_name, password )
VALUES ( ?, ?, ?, ? )
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
