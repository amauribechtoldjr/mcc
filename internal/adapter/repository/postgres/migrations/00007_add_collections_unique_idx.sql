-- +goose Up
ALTER TABLE collections 
ADD CONSTRAINT collection_idx UNIQUE ("name");

-- +goose Down
ALTER TABLE collections
DROP CONSTRAINT collection_idx;
