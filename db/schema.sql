CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(255) primary key);
CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE timings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    race_id integer,
    event_id integer,
    timing integer,
    recorded DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- Dbmate schema migrations
