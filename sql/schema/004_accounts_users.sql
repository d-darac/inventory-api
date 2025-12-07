-- +goose Up
CREATE TABLE accounts_users 
(
    account_id UUID NOT NULL REFERENCES accounts ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (account_id, user_id)
);

-- +goose Down
DROP TABLE accounts_users;