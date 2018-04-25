CREATE EXTENSION IF NOT EXISTS CITEXT;

drop table if exists users cascade;
drop table if exists forums cascade;
drop table if exists threads cascade;
drop table if exists posts cascade;
drop table if exists votes cascade;


CREATE TABLE IF NOT EXISTS users
(
  id       BIGSERIAL PRIMARY KEY,

  nickname VARCHAR(64) NOT NULL UNIQUE,
  email    CITEXT NOT NULL UNIQUE,

  about    TEXT DEFAULT '',
  fullname VARCHAR(96) DEFAULT ''
);


CREATE TABLE IF NOT EXISTS forums
(
  id      BIGSERIAL primary key,

  slug    CITEXT not null unique,

  title   CITEXT,

  threads INTEGER DEFAULT 0,
  posts   INTEGER DEFAULT 0,

  author  VARCHAR references users(nickname)
);

CREATE TABLE threads
(
  id         BIGSERIAL PRIMARY KEY,
  slug       CITEXT unique,

  created    TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,

  message    TEXT default '',
  title      TEXT default '',

  author     VARCHAR REFERENCES users (nickname),
  forum      CITEXT REFERENCES forums(slug),

  votes      BIGINT DEFAULT 0
);

create table if not exists posts
(
  id        bigserial not null primary key,

  created   TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,

  is_edited boolean default FALSE,

  parent    integer,
  path      bigint array,

  author    varchar not null references users(nickname),
  forum     CITEXT references forums(slug),
  thread    bigint references threads(id)
);

