-- migrate:up

CREATE TABLE races (
    race_id varchar(255) UNIQUE,
    owner_id varchar(255),
    race_name varchar(255)
);

CREATE TABLE events (
    race_id varchar(255),
    event_id varchar(255),
    event_description varchar(255),
    distance varchar(255),
    UNIQUE(race_id, event_id, event_description)
);

CREATE TABLE participants (
    race_id varchar(255) NOT NULL,
    event_id varchar(255) NOT NULL,
    bib_number varchar(255) NOT NULL,
    first_name varchar(255),
    last_name varchar(255),
    birth_year integer NOT NULL,
    gender varchar(255) NOT NULL,
    team varchar(255),
    UNIQUE(race_id, event_id, bib_number)
);

CREATE TABLE results (
    race_id varchar(255),
    event_id varchar(255),
    bib_number varchar(255),
    result integer,
    UNIQUE(race_id, event_id, bib_number)
);

-- migrate:down
DROP TABLE races;
DROP TABLE events;
DROP TABLE participants;
DROP TABLE results;
