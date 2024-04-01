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
VALUES ( ?, ?, ?, ?, ?, ?, ?, ? )
RETURNING *;

-- name: GetReviews :many
SELECT
    r.id,
    r.user_id,
    r.created_at,
    r.updated_at,
    u.display_name as user_name,
    r.movie_tmdb_id,
    r.rating,
    r.review 
FROM reviews r
JOIN users u
ON u.id = r.user_id
WHERE public_review = 1
LIMIT 5;
