-- name: GetUserAdminStats :one
SELECT 
    (SELECT COUNT(*) FROM movies) AS total_movies,
    (SELECT COUNT(*) FROM users) AS total_users,
    (SELECT COUNT(*) FROM bookings) AS total_bookings,
    CAST((SELECT COUNT(*) FROM booking_seats) * 10.99 AS DECIMAL(10, 2)) AS total_revenue;

-- name: GetTopFiveMovies :many
SELECT
    movies.id,
    movies.title,
    COUNT(bookings.id) AS total_bookings
FROM movies
LEFT JOIN showtimes ON movies.id = showtimes.movie_id
LEFT JOIN bookings ON showtimes.id = bookings.showtime_id
GROUP BY movies.id, movies.title
ORDER BY total_bookings DESC
LIMIT 5;

-- name: GetShowtimeOccupancy :many
SELECT 
    movies.title, 
    showtimes.start_time, 
    COUNT(booking_seats.seat_code) AS occupied_seats
FROM showtimes
LEFT JOIN movies ON showtimes.movie_id = movies.id
LEFT JOIN bookings ON showtimes.id = bookings.showtime_id
LEFT JOIN booking_seats ON bookings.id = booking_seats.booking_id
WHERE movies.title = $1
GROUP BY movies.title, showtimes.start_time
ORDER BY start_time ASC;
