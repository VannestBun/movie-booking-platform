-- name: CreateBooking :one
INSERT INTO bookings (
  id, 
  user_id, 
  showtime_id, 
  created_at, 
  updated_at
) VALUES (
  gen_random_uuid(), 
  $1, 
  $2,
  NOW(), 
  NOW()
)
RETURNING *;