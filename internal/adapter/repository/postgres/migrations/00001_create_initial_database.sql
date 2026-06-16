-- +goose Up
CREATE TABLE IF NOT EXISTS game (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  "name" TEXT NOT NULL,
  code TEXT NOT NULL
);

INSERT INTO game ("name", code)
VALUES 
  ('Magic: The Gathering', 'mtg'), 
  ('Pokémon Trading Card Game', 'ptcg');

CREATE TABLE IF NOT EXISTS "user" (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  "name" TEXT NOT NULL
);

INSERT INTO "user" ("name") VALUES ('Amauri Bechtold Junior');

CREATE TABLE IF NOT EXISTS "card" (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  oracle_id UUID NOT NULL,
  game_id UUID NOT NULL,
  CONSTRAINT card_unique_oracle_idx UNIQUE (oracle_id),
  CONSTRAINT fk_game_id FOREIGN KEY (game_id) REFERENCES game(id)
);

CREATE TABLE IF NOT EXISTS "collection" (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  "name" TEXT NOT NULL,
  game_id UUID NOT NULL,
  user_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT fk_game_id FOREIGN KEY (game_id) REFERENCES game(id),
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES "user"(id),
  CONSTRAINT collection_unique_name_idx UNIQUE ("name", user_id)
);

CREATE TABLE IF NOT EXISTS collection_card (
  card_id UUID NOT NULL,
  collection_id UUID NOT NULL,
  quantity SMALLINT NOT NULL,
  CONSTRAINT fk_card_id FOREIGN KEY (card_id) REFERENCES card(id),
  CONSTRAINT fk_collection_id FOREIGN KEY (collection_id) REFERENCES collection(id)
);

CREATE TABLE IF NOT EXISTS mtg_set (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  code TEXT NOT NULL,
  "name" TEXT NOT NULL,
  released_at TIMESTAMPTZ NOT NULL,
  parent_set_code TEXT,
  card_count INTEGER
);

CREATE TABLE IF NOT EXISTS mtg_card (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  set_id UUID NOT NULL,
  card_id UUID NOT NULL,
  name TEXT NOT NULL,
  layout TEXT,
  cmc NUMERIC(10,2),
  color_identity TEXT,
  color_indicator TEXT,
  colors TEXT,
  img_small_uri TEXT,
  img_normal_uri TEXT,
  CONSTRAINT fk_set_id FOREIGN KEY (set_id) REFERENCES mtg_set(id),
  CONSTRAINT fk_card_id FOREIGN KEY (card_id) REFERENCES "card"(id)
);

CREATE TABLE IF NOT EXISTS mtg_related (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  import_id UUID NOT NULL,
  layout TEXT NOT NULL,
  component TEXT NOT NULL,
  "name" TEXT NOT NULL,
  type_line TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS mtg_related_card (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  mtg_card_id UUID NOT NULL,
  mtg_related_id UUID NOT NULL,
  CONSTRAINT fk_mtg_card_id FOREIGN KEY (mtg_card_id) REFERENCES mtg_card(id),
  CONSTRAINT fk_mtg_related_id FOREIGN KEY (mtg_related_id) REFERENCES mtg_related(id)
);

-- +goose Down
DROP TABLE IF EXISTS "collection_card";
DROP TABLE IF EXISTS "mtg_related_card";
DROP TABLE IF EXISTS "mtg_related";
DROP TABLE IF EXISTS "mtg_card";
DROP TABLE IF EXISTS "mtg_set";
DROP TABLE IF EXISTS "collection";
DROP TABLE IF EXISTS "card";
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS "game";
