-- name: CreateLocation :one
INSERT INTO locations (name, notes) VALUES (?, ?) RETURNING id;

-- name: ListLocations :many
SELECT id, name, notes FROM locations;

-- name: DescribeLocation :one
SELECT name, notes FROM locations WHERE id=?;

-- name: UpdateLocationName :exec
UPDATE locations SET name=? WHERE id=?;

-- name: UpdateLocationNotes :exec
UPDATE locations SET notes=? WHERE id=?;

-- name: UpdateLocation :exec
UPDATE locations SET name=?, notes=? WHERE id=?;

-- name: DestroyLocation :exec
DELETE FROM locations WHERE id=?;
