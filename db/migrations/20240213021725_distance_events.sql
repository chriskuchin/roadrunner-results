-- migrate:up
CREATE TABLE attempts (
  race_id VARCHAR(255),
  event_id VARCHAR(255),
  bib VARCHAR(255),
  attempt_no INTEGER DEFAULT 1,
  result REAL,
  UNIQUE(race_id, event_id, bib, attempt_no)
);

-- migrate:down
drop table attempts;