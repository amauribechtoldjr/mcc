-- +goose Up
ALTER TABLE collections_cards 
ADD CONSTRAINT card_collection_idx UNIQUE (collection_id, card_id);

-- +goose Down
ALTER TABLE collections_cards
DROP CONSTRAINT card_collection_idx;
