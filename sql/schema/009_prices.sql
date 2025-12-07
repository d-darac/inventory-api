-- +goose Up
CREATE TABLE prices
( 
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    currency   currency NOT NULL,
    amount     INTEGER NOT NULL,
    account_id UUID NOT NULL REFERENCES accounts ON DELETE CASCADE,
    item_id    UUID NOT NULL REFERENCES items ON DELETE CASCADE
);

-- +goose Down
DROP TABLE prices;