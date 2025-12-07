-- name: CreateItemIdentifier :one
INSERT INTO item_identifiers 
(
    id,
    created_at,
    updated_at,
    ean,
    gtin,
    isbn,
    jan,
    mpn,
    nsn,
    upc,
    qr,
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
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
)
RETURNING *;