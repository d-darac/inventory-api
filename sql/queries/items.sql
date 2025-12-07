-- name: CreateItem :one
INSERT INTO items 
(
    id,
    created_at,
    updated_at,
    active,
    description,
    name,
    type,
    account_id,
    group_id
)
VALUES 
(
    gen_random_uuid(),
    NOW(),
    NOW(),
    TRUE,
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING
    id,
    created_at,
    updated_at,
    active,
    description,
    name,
    type,
    group_id AS group;
--

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = $1 AND account_id = $2;
--

-- name: GetItem :one
SELECT
    id,
    created_at,
    updated_at,
    active,
    description,
    name,
    type,
    group_id AS group
FROM items
WHERE id = $1 AND account_id =$2;
--

-- name: ListItems :many
SELECT
    id,
    created_at,
    updated_at,
    active,
    description,
    name,
    type,
    group_id AS group
FROM items
WHERE account_id = sqlc.arg('account_id')
AND
    (
        sqlc.narg('created_at_gt')::timestamp IS NULL 
        OR sqlc.narg('created_at_gt')::timestamp > created_at
    )
AND 
    (
        sqlc.narg('created_at_lt')::timestamp IS NULL 
        OR sqlc.narg('created_at_lt')::timestamp < created_at
    )
AND 
    (
        sqlc.narg('created_at_gte')::timestamp IS NULL 
        OR sqlc.narg('created_at_gte')::timestamp >= created_at
    )
AND 
    (
        sqlc.narg('created_at_lte')::timestamp IS NULL 
        OR sqlc.narg('created_at_lte')::timestamp <= created_at
    )
AND 
    (
        sqlc.narg('updated_at_gt')::timestamp IS NULL 
        OR sqlc.narg('updated_at_gt')::timestamp > updated_at
    )
AND 
    (
        sqlc.narg('updated_at_lt')::timestamp IS NULL 
        OR sqlc.narg('updated_at_lt')::timestamp < updated_at
    )
AND 
    (
        sqlc.narg('updated_at_gte')::timestamp IS NULL 
        OR sqlc.narg('updated_at_gte')::timestamp >= updated_at
    )
AND 
    (
        sqlc.narg('updated_at_lte')::timestamp IS NULL 
        OR sqlc.narg('updated_at_lte')::timestamp <= updated_at
    )
AND 
    (
        sqlc.narg('starting_after')::uuid IS NULL
        OR sqlc.narg('starting_after')::uuid > id
    )
AND 
    (
        sqlc.narg('ending_before')::uuid IS NULL
        OR sqlc.narg('ending_before')::uuid < id
    )
AND
    (
        sqlc.narg('active')::boolean IS NULL
        OR active = sqlc.narg('active')::boolean
    )
AND 
    (
        sqlc.narg('description')::text IS NULL 
        OR description LIKE sqlc.narg('description')::text
    )
AND 
    (
        sqlc.narg('group_id')::uuid IS NULL 
        OR group_id = sqlc.narg('group_id')::uuid
    )
AND 
    (
        sqlc.narg('name')::text IS NULL 
        OR name LIKE sqlc.narg('name')::text
    )
AND
    (
        sqlc.narg('type')::item_type IS NULL 
        OR name LIKE sqlc.narg('type')::item_type
    )
ORDER BY created_at DESC
LIMIT COALESCE(sqlc.narg('limit')::int, 10);
--

-- name: UpdateItem :one
UPDATE items
SET
    updated_at = NOW(),
    active = COALESCE(sqlc.narg('active'), active),
    description = COALESCE(sqlc.narg('description'), description),
    name = COALESCE(sqlc.narg('name'), name),
    group_id = COALESCE(sqlc.narg('group_id'), group_id)
WHERE id = sqlc.arg('id') AND account_id = sqlc.arg('account_id')
RETURNING
    id,
    created_at,
    updated_at,
    active,
    description,
    name,
    type,
    group_id AS group;
--
