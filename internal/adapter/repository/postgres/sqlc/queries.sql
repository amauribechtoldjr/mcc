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

-- name: UpsertMTGCard :batchexec
WITH existing AS (
  SELECT id FROM mtg_card
  WHERE set_id = $2 AND collector_number = $5 AND lang = $4
),
new_card AS (
  INSERT INTO "card" (game_id)
  SELECT $1
  WHERE NOT EXISTS (SELECT 1 FROM existing)
  RETURNING id
),
resolved AS (
  SELECT id FROM existing
  UNION ALL
  SELECT id FROM new_card
)
INSERT INTO mtg_card (
  id,
  set_id,
  oracle_id,
  lang,
  collector_number,
  "name",
  printed_type_line,
  printed_text,
  flavor_text,
  layout,
  cmc,
  color_identity,
  color_indicator,
  colors,
  img_small_uri,
  img_normal_uri,
  last_import_id
)
SELECT id, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
FROM resolved
ON CONFLICT (set_id, collector_number, lang) DO UPDATE SET
  oracle_id = EXCLUDED.oracle_id,
  "name" = EXCLUDED.name,
  printed_type_line = EXCLUDED.printed_type_line,
  printed_text = EXCLUDED.printed_text,
  flavor_text = EXCLUDED.flavor_text,
  layout = EXCLUDED.layout,
  cmc = EXCLUDED.cmc,
  color_identity = EXCLUDED.color_identity,
  color_indicator = EXCLUDED.color_indicator,
  colors = EXCLUDED.colors,
  img_small_uri = EXCLUDED.img_small_uri,
  img_normal_uri = EXCLUDED.img_normal_uri,
  last_import_id = EXCLUDED.last_import_id;

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
ON CONFLICT (code) DO UPDATE SET
  import_id = EXCLUDED.import_id,
  "name" = EXCLUDED.name,
  released_at = EXCLUDED.released_at,
  parent_set_code = EXCLUDED.parent_set_code,
  card_count = EXCLUDED.card_count
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