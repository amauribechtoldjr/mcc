-- name: ListCards :many
SELECT * FROM cards;

-- name: FindCardById :one
SELECT * from cards WHERE id = $1;