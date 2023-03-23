-- migrate:up
CREATE TABLE timers (
    timer_id varchar(255),
    race_id varchar(255),
    event_id varchar(255),
    start_ts integer
);

CREATE TABLE timer_results (
    timer_id varchar(255),
    race_id varchar(255),
    event_id varchar(255),
    result integer
);

-- migrate:down
DROP TABLE timers;
DROP TABLE timer_results

