CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE races (
    race_id varchar(255) UNIQUE,
    owner_id varchar(255),
    race_name varchar(255)
);
CREATE TABLE events (
    race_id varchar(255),
    event_id varchar(255),
    event_description varchar(255),
    event_type varchar(255),
    distance varchar(255) DEFAULT NULL,
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
    timer_id varchar(255) DEFAULT NULL,
    bib_number varchar(255),
    result integer,
    UNIQUE(race_id, event_id, bib_number)
);
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
CREATE TABLE users (
    user_id varchar(255),
    email varchar(255),
    display_name TEXT,
    creation_date TEXT DEFAULT CURRENT_TIMESTAMP,
    last_login_date TEXT
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20220708034230'),
  ('20230321020148'),
  ('20230813205300');
