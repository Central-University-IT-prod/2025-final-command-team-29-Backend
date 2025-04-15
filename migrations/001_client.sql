-- +goose Up
CREATE TABLE client (
    id uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    UNIQUE (name),
    UNIQUE (phone_number)
);