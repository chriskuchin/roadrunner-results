CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(255) primary key);
CREATE TABLE sqlite_sequence(name,seq);
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
    bib_number varchar(255) NOT NULL,
    first_name varchar(255),
    last_name varchar(255),
    birth_year integer NOT NULL,
    gender varchar(255) NOT NULL,
    team varchar(255),
    UNIQUE(race_id, bib_number)
);
CREATE TABLE event_participation (
    race_id varchar(255) NOT NULL,
    event_id varchar(255) NOT NULL,
    bib_number varchar(255) NOT NULL,
    UNIQUE(race_id, event_id, bib_number)
);
CREATE TABLE results (
    race_id varchar(255),
    event_id varchar(255),
    bib_number varchar(255),
    result integer,
    UNIQUE(race_id, event_id, bib_number)
);
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
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20220708034230'),
  ('20230321020148');
