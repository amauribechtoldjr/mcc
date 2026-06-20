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

-- name: FindGameByCode :one
SELECT * FROM game WHERE code = $1;

-- name: CreateCard :one
INSERT INTO "card" (oracle_id, game_id)
VALUES ($1, $2)
ON CONFLICT (oracle_id) DO UPDATE
  SET oracle_id = EXCLUDED.oracle_id
RETURNING id;

-- name: CreateMTGCard :exec
INSERT INTO mtg_card (
  set_id, 
  card_id,
  "name",
  layout, 
  cmc, 
  color_identity, 
  color_indicator, 
  colors, 
  img_small_uri, 
  img_normal_uri,
  last_import_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: CreateMTGSet :one
INSERT INTO mtg_set (
  import_id,
  "name",
  code,
  released_at,
  parent_set_code, 
  card_count
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: CreateScryfallImport :one
INSERT INTO scryfall_import (
  started_at,
  bulk_updated_at,
  "status"
)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetScryfallImportCount :one
SELECT 
  count(id) as import_quantity 
FROM 
  scryfall_import 
WHERE 
  bulk_updated_at >= $1;

-- name: UpdateScryfallImport :exec
UPDATE
  scryfall_import
SET
  finished_at = $1, "status" = $2
WHERE
  id = $3;