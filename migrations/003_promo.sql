-- +goose Up
CREATE TABLE promos (
    id uuid PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    partner_id uuid,
    condition_string VARCHAR(255),
    condition JSONB, 
    usage_limit INTEGER,
    time_limit_start BIGINT,
    time_limit_end BIGINT,
    is_active BOOLEAN,
    is_approved BOOLEAN,
    image_id INTEGER NOT NULL,
    FOREIGN KEY (partner_id) REFERENCES partner(id)
);