-- +goose Up
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  name VARCHAR(100) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
