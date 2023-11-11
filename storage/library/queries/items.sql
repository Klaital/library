-- name: ListItemsForLocation :many
SELECT id, code, code_type, code_source,
       title, title_translated, title_transliterated,
       created_at, updated_at
FROM items
WHERE location_id=?;

-- name: CreateItem :one
INSERT INTO items (location_id, code, code_type, code_source,
                   title, title_translated, title_transliterated)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: ListAllItems :many
SELECT id, code, code_type, code_source,
       title, title_translated, title_transliterated,
       created_at, updated_at
FROM items;

-- name: GetItem :one
SELECT id, code, code_type, code_source,
       title, title_translated, title_transliterated,
       created_at, updated_at
FROM items
WHERE id=?;

-- name: MoveItem :exec
UPDATE items SET location_id = ? WHERE id = ? LIMIT 1;
