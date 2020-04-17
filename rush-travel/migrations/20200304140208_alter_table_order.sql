  
-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE orders
ADD COLUMN status varchar;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE orders
DROP COLUMN status;