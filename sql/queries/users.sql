-- name: CreateUser :one
INSERT INTO users 
(
    id,
    created_at,
    updated_at,
    email,
    hashed_password,
    name
)
VALUES 
(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING
    id,
    created_at,
    updated_at,
    email,
    name;
--

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
--

-- name: GetUser :one
SELECT
    id,
    created_at,
    updated_at,
    email,
    name
FROM users
WHERE id = $1;
--

-- name: GetUserByEmail :one
SELECT
    id,
    created_at,
    updated_at,
    email,
    name
FROM users
WHERE email = $1;
--

-- name: UpdateUser :one
UPDATE users
SET
    updated_at = NOW(),
    email = COALESCE(sqlc.narg('email'), email),
    hashed_password = COALESCE(sqlc.narg('hashed_password'), hashed_password),
    name = COALESCE(sqlc.narg('name'), name)
WHERE id = sqlc.arg('id')
RETURNING
    id,
    created_at,
    updated_at,
    email,
    name;
--