-- +goose Up
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS bookings_movie_id_seat_number_booking_time_key;
ALTER TABLE bookings DROP COLUMN movie_id;
ALTER TABLE bookings DROP COLUMN seat_number;
ALTER TABLE bookings DROP COLUMN booking_time;
ALTER TABLE bookings ADD COLUMN showtime_id UUID REFERENCES showtimes(id);

-- +goose Down
ALTER TABLE bookings ADD COLUMN movie_id UUID REFERENCES movies(id);
ALTER TABLE bookings ADD COLUMN seat_number VARCHAR(10);
ALTER TABLE bookings ADD COLUMN booking_time TIMESTAMP NOT NULL;
ALTER TABLE bookings DROP COLUMN showtime_id;
ALTER TABLE bookings ADD CONSTRAINT bookings_movie_id_seat_number_booking_time_key UNIQUE(movie_id, seat_number, booking_time);