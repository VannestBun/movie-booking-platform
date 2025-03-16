-- +goose Up
ALTER TABLE booking_seats ALTER COLUMN booking_id SET NOT NULL;

-- +goose Down
ALTER TABLE booking_seats ALTER COLUMN booking_id SET NULL;