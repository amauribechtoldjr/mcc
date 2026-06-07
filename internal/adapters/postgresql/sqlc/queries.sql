-- name: ListCards :many
SELECT * FROM cards;

-- name: FindCardById :one
SELECT * from cards WHERE id = $1;

-- name: CreateCollection :one
INSERT INTO collections (user_id, "name") VALUES ($1, $2) RETURNING *;

-- name: AddCardToCollection :exec
INSERT INTO collections_cards (card_id, collection_id, quantity) VALUES ($1, $2, $3);