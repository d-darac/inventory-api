-- +goose Up
CREATE TABLE inventories 
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    in_stock   INTEGER NOT NULL DEFAULT 0,
    orderable  INTEGER,
    account_id UUID NOT NULL REFERENCES accounts ON DELETE CASCADE,
    item_id    UUID NOT NULL REFERENCES items ON DELETE CASCADE
);

-- +goose Down
DROP TABLE inventories;