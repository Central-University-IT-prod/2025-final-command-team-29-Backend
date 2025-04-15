-- +goose Up
CREATE TABLE activations (
    client_id uuid NOT NULL,
    promo_id uuid NOT NULL,
    buy INT, 
    current INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL,
    overall INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (client_id) REFERENCES client(id),
    FOREIGN KEY (promo_id) REFERENCES promos(id) ON DELETE CASCADE
);
