-- name: ListCards :many
SELECT * FROM "card";

-- name: FindCardById :one
SELECT * from "card" WHERE id = $1;

-- name: ListCollections :many
SELECT * FROM "collection" WHERE user_id = $1;

-- name: CreateCollection :one
INSERT INTO "collection" (user_id, "name") VALUES ($1, $2) RETURNING *;

-- name: AddCardToCollection :exec
INSERT INTO collection_card (card_id, collection_id, quantity) VALUES ($1, $2, $3);