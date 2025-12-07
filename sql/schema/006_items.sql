-- +goose Up
CREATE TABLE items 
(
    id             UUID PRIMARY KEY,
    created_at     TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP NOT NULL,
    active         BOOLEAN NOT NULL DEFAULT TRUE,
    description    TEXT,
    name           TEXT NOT NULL,
    type           item_type NOT NULL,
    account_id     UUID NOT NULL REFERENCES accounts ON DELETE CASCADE,
    group_id       UUID DEFAULT NULL REFERENCES groups ON DELETE SET NULL
);

-- +goose Down
DROP TABLE items;