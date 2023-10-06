-- migrate:up
CREATE TABLE race_divisions (
  race_id VARCHAR(255),
  display VARCHAR(255),
  filters TEXT
);

-- migrate:down
DROP TABLE race_divisions;