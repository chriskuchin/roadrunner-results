-- migrate:up
CREATE TABLE participants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    race_id integer,
    first_name varchar(255),
    last_name varchar(255),
    team varchar(255),
    birth_year integer
);

CREATE TABLE races (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    race_id varchar(255),
    owner_id varchar(255),
    race_name varchar(255),
);

CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_description varchar(255),
    distance varchar(255)
);

CREATE TABLE results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    race_id integer,
    event_id integer,
    bib_number integer,
    recorded DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE timings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    race_id integer,
    event_id integer,
    timing integer,
    recorded DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id varchar(255),
    email varchar(255)
);

-- migrate:down
DROP TABLE participants;
DROP TABLE races;
DROP TABLE events;
DROP TABLE results;
DROP TABLE timings;
DROP TABLE users;