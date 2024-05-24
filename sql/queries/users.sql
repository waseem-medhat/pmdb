-- name: GetUser :one
SELECT id, display_name, user_name, bio FROM users
WHERE user_name = $1 LIMIT 1;

-- name: GetUserForLogin :one
SELECT user_name, password FROM users
WHERE user_name = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY display_name;

-- name: CreateUser :one
INSERT INTO users ( id, user_name, display_name, password )
VALUES ( $1, $2, $3, $4 )
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
