-- +goose Up
ALTER TABLE users
ADD COLUMN user_role TEXT NOT NULL DEFAULT 'user';

-- +goose Down
ALTER TABLE users
DROP COLUMN user_role;