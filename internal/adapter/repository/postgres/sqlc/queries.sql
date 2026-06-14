-- name: ListCards :many
SELECT * FROM "card";

-- name: FindCardById :one
SELECT * from "card" WHERE id = $1;

-- name: ListCollections :many
SELECT * FROM "collection" WHERE user_id = $1;

-- name: CreateCollection :one
INSERT INTO "collection" (user_id, "name") VALUES ($1, $2) RETURNING *;

-- name: AddCardToCollection :exec
INSERT INTO collections_card (card_id, collection_id, quantity) VALUES ($1, $2, $3);

-- name: ListCollectionCards :many
SELECT
	c.id,
  c."name",
  cc.quantity
FROM
	collections_card cc
INNER JOIN cards c ON c.id = cc.card_id
WHERE
  cc.collection_id = $1;