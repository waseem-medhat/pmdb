-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY display_name;

-- name: CreateUser :one
INSERT INTO users ( id, user_name, display_name )
VALUES ( ?, ?, ? )
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
