-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users
ALTER COLUMN phone_number TYPE varchar;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users
ALTER COLUMN phone_number TYPE int;