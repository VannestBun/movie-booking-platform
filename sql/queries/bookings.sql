-- name: CreateBooking :one
INSERT INTO bookings (
  id, 
  user_id, 
  movie_id, 
  seat_number, 
  booking_time, 
  created_at, 
  updated_at
) VALUES (
  gen_random_uuid(), 
  $1, 
  $2, 
  $3, 
  $4, 
  NOW(), 
  NOW()
)
RETURNING *;