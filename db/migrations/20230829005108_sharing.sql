-- migrate:up
CREATE TABLE race_authorization (
  user_id VARCHAR(255),
  race_id VARCHAR(255),
  permissions TEXT
);

-- migrate:down
DROP TABLE race_authz;