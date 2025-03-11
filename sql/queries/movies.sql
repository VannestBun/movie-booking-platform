-- name: CreateMovie :one
INSERT INTO movies (
    id,
    created_at,
    updated_at,
    title,
    description,
    duration_minutes,
    poster_image_url,
    trailer_video_url
) VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,  -- title
    $2,  -- description
    $3,  -- duration_minutes
    $4,  -- poster_image_url
    $5   -- trailer_video_url
)
RETURNING *;

-- name: GetMovie :one
Select *
FROM movies
WHERE id = $1;

-- name: GetMovies :many
SELECT * FROM movies
ORDER BY title ASC;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1;

-- name: UpdateMovie :exec
UPDATE movies
SET title = $1, description = $2, duration_minutes = $3, poster_image_url = $4, trailer_video_url = $5
WHERE id = $6;