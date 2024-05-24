-- name: CreateReview :one
INSERT INTO reviews (
    id,
    created_at,
    updated_at,
    user_id,
    movie_tmdb_id,
    rating,
    review,
    public_review
)
VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 )
RETURNING *;

-- name: GetReviews :many
SELECT
    r.id,
    r.user_id,
    u.display_name as user_name,
    r.created_at,
    r.updated_at,
    r.movie_tmdb_id,
    r.rating,
    r.review 
FROM reviews r
JOIN users u
ON u.id = r.user_id
WHERE public_review
ORDER BY r.created_at DESC
LIMIT 5;

-- name: GetReviewsForMovie :many
SELECT
    r.id,
    r.user_id,
    u.display_name as user_name,
    r.created_at,
    r.updated_at,
    r.movie_tmdb_id,
    r.rating,
    r.review 
FROM reviews r
JOIN users u
ON u.id = r.user_id
WHERE movie_tmdb_id = $1 AND public_review;

-- name: GetReviewByID :one
SELECT
    r.id,
    r.user_id,
    u.display_name as user_name,
    r.created_at,
    r.updated_at,
    r.movie_tmdb_id,
    r.rating,
    r.review 
FROM reviews r
JOIN users u
ON u.id = r.user_id
WHERE r.id = $1;

