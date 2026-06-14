-- +goose Up
ALTER TABLE "collection"
ADD created_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- +goose Down
ALTER TABLE "collection"
DROP created_at;
