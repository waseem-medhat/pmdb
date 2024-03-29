-- name: CreateReview :one
INSERT INTO reviews ( id, user_id, movie_tmdb_id, rating, review, public_review )
VALUES ( ?, ?, ?, ?, ?, ? )
RETURNING *;

-- name: GetReviews :many
SELECT  id, user_id, movie_tmdb_id, rating, review 
FROM reviews
WHERE public_review = 1
LIMIT 5;
