-- +goose Up
ALTER TABLE mtg_set
ADD COLUMN import_id UUID NOT NULL;

-- +goose Down
ALTER TABLE mtg_set
DROP COLUMN import_id;
