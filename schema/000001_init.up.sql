
CREATE TABLE users (
    guid varchar(200) not null unique,
    username varchar(200) not null unique,
    email varchar(200) not null unique,
    password_hash varchar(256) not null
);


CREATE TABLE refresh_tokens (
    user_id varchar(200) references users (guid) on delete cascade not null,
    token varchar(256) not null unique
);