-- migrate:up
ALTER TABLE participants ADD COLUMN grade INTEGER DEFAULT NULL;

-- migrate:down
ALTER TABLE participants REMOVE COLUMN grade;
