-- +goose Up
CREATE TABLE movies (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    duration_minutes INTEGER NOT NULL,
    poster_image_url TEXT NOT NULL,
    trailer_video_url TEXT NOT NULL
);

-- +goose Down
DROP TABLE movies;