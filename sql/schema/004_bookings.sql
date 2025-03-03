-- +goose Up
CREATE TABLE bookings (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  movie_id UUID NOT NULL REFERENCES movies(id),
  booking_time TIMESTAMP NOT NULL,
  seat_number VARCHAR(10) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  UNIQUE(movie_id, seat_number, booking_time) -- ensure no double bookings
);

-- +goose Down
DROP TABLE bookings;