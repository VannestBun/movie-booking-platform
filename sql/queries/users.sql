-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUser :one
Select *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
Select *
FROM users
WHERE email = $1;

-- name: GetUserBookingInfo :many
SELECT 
    users.email,
    showtimes.start_time, 
    movies.title, 
    movies.description,
    movies.duration_minutes,
    movies.poster_image_url,
    movies.trailer_video_url,
    movies.rating,
    movies.genre,
    movies.director,
    movies.casts,
    bookings.id as booking_id
FROM users
LEFT JOIN bookings ON users.id = bookings.user_id
LEFT JOIN showtimes ON bookings.showtime_id = showtimes.id
LEFT JOIN movies ON showtimes.movie_id = movies.id
WHERE users.id = $1;
