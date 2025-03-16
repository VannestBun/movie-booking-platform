-- +goose Up
ALTER TABLE bookings ALTER COLUMN showtime_id SET NOT NULL;

-- +goose Down
ALTER TABLE bookings ALTER COLUMN showtime_id SET NULL;