-- name: CreateBookingSeat :one
INSERT INTO booking_seats (id, booking_id, seat_code)
VALUES (
    gen_random_uuid(),
    $1,
    $2
)
RETURNING *;

-- name: GetBookedSeats :many
SELECT seat_code
FROM booking_seats
LEFT JOIN bookings ON booking_seats.booking_id = bookings.id
LEFT JOIN showtimes ON bookings.showtime_id = showtimes.id
WHERE showtimes.start_time = $1 AND showtimes.movie_id = $2;