-- +goose Up
CREATE TABLE showtimes (
    id UUID PRIMARY KEY,
    movie_id UUID REFERENCES movies(id),
    start_time TIMESTAMP NOT NULL,

    UNIQUE(movie_id, start_time)
);

-- +goose Down
DROP TABLE showtimes;