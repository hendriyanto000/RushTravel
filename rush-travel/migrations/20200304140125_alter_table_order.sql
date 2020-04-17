-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE orders
ALTER COLUMN price TYPE numeric;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE orders
ALTER COLUMN price TYPE int;