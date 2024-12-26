-- +goose Up
ALTER TABLE users ADD COLUMN is_chirpy_red BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE users DROP COLUMN is_chirpy_red;