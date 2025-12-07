-- name: CreateAccountUserReference :exec
INSERT INTO accounts_users 
(
    account_id,
    user_id,
    created_at,
    updated_at
)
VALUES 
(
    $1,
    $2,
    NOW(),
    NOW()
);