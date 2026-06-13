-- +goose Up
ALTER TABLE collections 
ADD COLUMN user_id UUID NOT NULL;

ALTER TABLE collections 
ADD CONSTRAINT fk_user 
FOREIGN KEY (user_id) 
REFERENCES users(id) 
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- +goose Down
ALTER TABLE collections
DROP CONSTRAINT fk_user;

ALTER TABLE collections
DROP COLUMN user_id;
