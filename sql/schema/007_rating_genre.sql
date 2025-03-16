-- +goose Up
ALTER TABLE movies
ADD COLUMN rating DECIMAL(2, 1) NOT NULL;

ALTER TABLE movies
ADD COLUMN genre TEXT NOT NULL;

-- +goose Down
ALTER TABLE movies
DROP COLUMN rating;

ALTER TABLE movies
DROP COLUMN genre;

