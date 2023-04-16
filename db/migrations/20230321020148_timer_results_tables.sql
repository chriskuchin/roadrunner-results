-- migrate:up
CREATE TABLE heats (
    race_id varchar(255),
    event_id varchar(255),
    timer_id varchar(255),
    heat varchar(255),
    bib_number varchar(255),
    lane_number varchar(255)
);

CREATE TABLE timers (
    race_id varchar(255),
    event_id varchar(255),
    timer_id varchar(255),
    start_ts integer DEFAULT 0
);

CREATE TABLE timer_results (
    race_id varchar(255),
    event_id varchar(255),
    timer_id varchar(255),
    result integer
);

-- migrate:down
DROP TABLE timers;
DROP TABLE timer_results;
DROP TABLE heats;