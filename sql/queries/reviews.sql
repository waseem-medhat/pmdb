-- name: CreateReview :one
INSERT INTO reviews ( id, user_id, movie_tmdb_id, rating, review, public_review )
VALUES ( ?, ?, ?, ?, ?, ? )
RETURNING *;
