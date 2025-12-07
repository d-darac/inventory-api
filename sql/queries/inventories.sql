-- name: CreateInventory :one
INSERT INTO inventories 
(
    id,
    created_at,
    updated_at,
    in_stock,
    orderable,
    account_id,
    item_id
)
VALUES 
(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
)
RETURNING
    id,
    created_at,
    updated_at,
    in_stock,
    orderable,
    item_id AS item;
--

-- name: DeleteInventory :exec
DELETE FROM inventories
WHERE id = $1 AND account_id = $2;
--

-- name: GetInventory :one
SELECT
    id,
    created_at,
    updated_at,
    in_stock,
    orderable,
    item_id AS item
FROM inventories
WHERE id = $1 AND account_id = $2;
--

-- name: ListInventories :many
SELECT
    id,
    created_at,
    updated_at,
    in_stock,
    orderable,
    item_id AS item
FROM inventories
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
ORDER BY created_at DESC
LIMIT COALESCE(sqlc.narg('limit')::int, 10);
--

-- name: UpdateInventory :one
UPDATE inventories
SET
    updated_at = NOW(),
    in_stock = COALESCE(sqlc.narg('in_stock'), in_stock),
    orderable = COALESCE(sqlc.narg('orderable'), orderable)
WHERE id = sqlc.arg('id') AND account_id = sqlc.arg('account_id')
RETURNING
    id,
    created_at,
    updated_at,
    in_stock,
    orderable,
    item_id AS item;
--
