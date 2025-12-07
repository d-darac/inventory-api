-- name: CreateGroup :one
INSERT INTO groups 
(
    id,
    created_at,
    updated_at,
    description,
    name,
    parent_id,
    account_id
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
    description,
    name,
    parent_id AS parent_group;
--

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1 AND account_id = $2;
--

-- name: GetGroup :one
SELECT
    id,
    created_at,
    updated_at,
    description,
    name,
    parent_id AS parent_group
FROM groups
WHERE id = $1 AND account_id = $2;
--

-- name: ListGroups :many
SELECT
    id,
    created_at,
    updated_at,
    description,
    name,
    parent_id AS parent_group
FROM groups
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
        sqlc.narg('parent_id')::uuid IS NULL 
        OR parent_id = sqlc.narg('parent_id')::uuid
    )
AND 
    (
        sqlc.narg('description')::text IS NULL 
        OR description LIKE sqlc.narg('description')::text
    )
AND 
    (
        sqlc.narg('name')::text IS NULL 
        OR name LIKE sqlc.narg('name')::text
    )
ORDER BY created_at DESC
LIMIT COALESCE(sqlc.narg('limit')::int, 10);
--

-- name: UpdateGroup :one
UPDATE groups
SET
    updated_at = NOW(),
    description = COALESCE(sqlc.narg('description'), description),
    name = COALESCE(sqlc.narg('name'), name),
    parent_id = COALESCE(sqlc.narg('parent_id'), parent_id)
WHERE id = sqlc.arg('id') AND account_id = sqlc.arg('account_id')
RETURNING
    id,
    created_at,
    updated_at,
    description,
    name,
    parent_id AS parent_group;
--
