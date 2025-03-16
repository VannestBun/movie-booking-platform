-- name: CreateShowtime :one
INSERT INTO showtimes (id, movie_id, start_time)
VALUES (
    gen_random_uuid(),
    $1,
    $2
)
RETURNING *;

-- name: GetShowtime :one
SELECT * FROM showtimes WHERE id = $1;

-- name: GetShowtimeByMovieID :one
SELECT * FROM showtimes WHERE movie_id = $1;

-- name: GetShowtimeByMovieAndStartTime :one
SELECT * FROM showtimes WHERE movie_id = $1 AND start_time = $2;