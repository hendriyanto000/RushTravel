-- +goose Up
CREATE TABLE Loggings(
 token varchar,
 username varchar,
 userStatus boolean,
 created_at timestamp,
 deleted_at timestamp 
);
-- +goose Down
DROP TABLE Loggings;