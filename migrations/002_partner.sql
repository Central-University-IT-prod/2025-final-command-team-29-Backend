-- +goose Up
CREATE TABLE partner (
    id uuid PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(500) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    image_id INTEGER NOT NULL,
    UNIQUE (title),
    UNIQUE (phone_number)
);