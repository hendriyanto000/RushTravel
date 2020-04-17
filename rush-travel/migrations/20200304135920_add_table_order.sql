-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE orders(
    id serial NOT NULL PRIMARY KEY,
    order_id varchar NOT NULL,
    username varchar,
    pick_up_loc text[] NOT NULL,
    drop_off_loc text[] NOT NULL,
    date DATE NOT NULL,
    days text NOT NULL,
    no_off_passangers int NOT NULL,
    price int NOT NULL,
    vehicle_type varchar
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE orders;
