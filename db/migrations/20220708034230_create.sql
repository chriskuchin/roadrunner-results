-- migrate:up
CREATE TABLE participants {
    id integer,
    raceID integer,
    firstName varchar(255),
    lastName varchar(255),
    team varchar(255),
    birthYear integer,
}

CREATE TABLE races {
    id integer,
    race varchar(255),
}

CREATE TABLE events {
    id integer,
    event varchar(255),
    measurement varchar(255),
}

CREATE TABLE results {
    id integer,
    raceID integer,
    participantID integer,
    eventID integer,
}

-- migrate:down
DROP TABLE participants;
DROP TABLE races;
DROP TABLE events;
DROP TABLE results;