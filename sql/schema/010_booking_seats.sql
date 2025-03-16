-- +goose Up
CREATE TABLE booking_seats (
    id UUID PRIMARY KEY,
    booking_id UUID REFERENCES bookings(id),
    seat_code TEXT NOT NULL, -- e.g., "A1", "B5"
    UNIQUE(booking_id, seat_code)
);

-- +goose Down
DROP TABLE booking_seats;