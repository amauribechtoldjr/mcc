-- +goose Up
CREATE TABLE IF NOT EXISTS collection_cards (
  card_id UUID NOT NULL,
  collection_id UUID NOT NULL,
  quantity SMALLINT NOT NULL,
  CONSTRAINT fk_card_id FOREIGN KEY (card_id) REFERENCES cards(id),
  CONSTRAINT fk_collection_id FOREIGN KEY (collection_id) REFERENCES collections(id)
);

-- +goose Down
DROP TABLE IF EXISTS collection_cards;