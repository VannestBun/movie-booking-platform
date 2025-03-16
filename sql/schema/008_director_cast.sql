-- +goose Up
ALTER TABLE movies
ADD COLUMN director TEXT NOT NULL;

ALTER TABLE movies
ADD COLUMN casts TEXT[] NOT NULL;

-- +goose Down
ALTER TABLE movies
DROP COLUMN director;

ALTER TABLE movies
DROP COLUMN casts;

