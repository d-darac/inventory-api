-- +goose Up
CREATE TABLE groups
(
    id          UUID PRIMARY KEY,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,
    description TEXT,
    name        TEXT NOT NULL,
    parent_id   UUID REFERENCES groups,
    account_id  UUID NOT NULL REFERENCES accounts ON DELETE CASCADE
);

-- +goose Down
DROP TABLE groups;