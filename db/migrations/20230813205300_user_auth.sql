-- migrate:up
CREATE TABLE users (
    user_id varchar(255),
    email varchar(255),
    display_name TEXT,
    creation_date TEXT DEFAULT CURRENT_TIMESTAMP,
    last_login_date TEXT
);

-- migrate:down
drop table users;