-- name: CreateReview :one
INSERT INTO reviews ( id, user_id, movie_tmdb_id, rating, review, public_review )
VALUES ( ?, ?, ?, ?, ?, ? )
RETURNING *;

-- name: GetReviews :many
SELECT
    r.id,
    r.user_id,
    u.display_name as user_name,
    r.movie_tmdb_id,
    r.rating,
    r.review 
FROM reviews r
JOIN users u
ON u.id = r.user_id
WHERE public_review = 1
LIMIT 5;
