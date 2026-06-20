-- +goose Up
CREATE TABLE IF NOT EXISTS scryfall_import (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  finished_at TIMESTAMPTZ,
  bulk_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  "status" TEXT NOT NULL
);

ALTER TABLE mtg_card
ADD COLUMN last_import_id UUID NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS scryfall_import;
ALTER TABLE mtg_card
DROP COLUMN last_import_id;
