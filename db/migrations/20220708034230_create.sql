-- migrate:up
CREATE TABLE members (
    first_name varchar(255),
    last_name varchar(255),
    birth_year integer,
    gender varchar(255),
    UNIQUE(first_name, last_name, birth_year)
);

CREATE TABLE races (
    race_id varchar(255) UNIQUE,
    owner_id varchar(255),
    race_name varchar(255)
);

CREATE TABLE events (
    race_id varchar(255),
    event_description varchar(255),
    distance varchar(255)
);

CREATE TABLE participants (
    race_id varchar(255),
    event_id varchar(255),
    member_id varchar(255),
    team varchar(255),
    UNIQUE(race_id, event_id, member_id)
);

CREATE TABLE results (
    race_id integer,
    event_id integer,
    bib_number integer,
    result integer,
    UNIQUE(race_id, event_id, bib_number)
);

CREATE TABLE timings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    race_id integer,
    event_id integer,
    timing integer,
    recorded DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE members;
DROP TABLE participants;
DROP TABLE races;
DROP TABLE events;
DROP TABLE results;
DROP TABLE timings;
