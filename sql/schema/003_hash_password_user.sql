-- +goose Up
ALTER TABLE users ADD COLUMN hashed_password TEXT;
UPDATE users SET hashed_password = 'unset';
ALTER TABLE users ALTER COLUMN hashed_password SET NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN hashed_password;