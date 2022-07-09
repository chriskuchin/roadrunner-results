-- migrate:up
CREATE TABLE participants (
    first_name varchar(255),
    last_name varchar(255),
    team varchar(255),
    birth_year integer
);

CREATE TABLE races (
    race varchar(255)
);

CREATE TABLE events (
    event_description varchar(255),
    measurement varchar(255)
);

CREATE TABLE results (
    race_id integer,
    participant_id integer,
    event_id integer
);

-- migrate:down
DROP TABLE participants;
DROP TABLE races;
DROP TABLE events;
DROP TABLE results;