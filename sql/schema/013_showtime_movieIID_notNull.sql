-- +goose Up
ALTER TABLE showtimes ALTER COLUMN movie_id SET NOT NULL;

-- +goose Down
ALTER TABLE showtimes ALTER COLUMN movie_id SET NULL;