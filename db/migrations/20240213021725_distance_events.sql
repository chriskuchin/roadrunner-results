-- migrate:up
CREATE TABLE attempts (
  race_id VARCHAR(255),
  event_id VARCHAR(255),
  bib VARCHAR(255),
  result real,
);

-- migrate:down
drop table attempts;