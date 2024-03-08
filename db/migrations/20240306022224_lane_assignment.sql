-- migrate:up
CREATE TABLE lane_assignments (
  race_id varchar(255),
  event_id varchar(255),
  timer_id varchar(255),
  assignments TEXT,
  UNIQUE(race_id, event_id, timer_id)
);

-- migrate:down
drop table lane_assignments;