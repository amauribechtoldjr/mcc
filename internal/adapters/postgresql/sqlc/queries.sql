-- name: ListCards :many
SELECT * FROM cards;

-- name: FindCardById :one
SELECT * from cards WHERE id = $1;

-- name: CreateCollection :one
INSERT INTO collections (user_id, "name") VALUES ($1, $2) RETURNING *;