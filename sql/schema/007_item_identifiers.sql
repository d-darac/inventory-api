-- +goose Up
CREATE TABLE item_identifiers
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    ean        TEXT,
    gtin       TEXT,
    isbn       TEXT,
    jan        TEXT,
    mpn        TEXT,
    nsn        TEXT,
    upc        TEXT,
    qr         TEXT,
    account_id UUID NOT NULL REFERENCES accounts ON DELETE CASCADE,
    item_id    UUID NOT NULL REFERENCES items ON DELETE CASCADE
);

-- +goose Down
DROP TABLE item_identifiers;