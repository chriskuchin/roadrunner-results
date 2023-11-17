-- migrate:up
ALTER TABLE races ADD COLUMN race_date INTEGER DEFAULT NULL;

-- migrate:down
ALTER TABLE races REMOVE COLUMN race_date;
