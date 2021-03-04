-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table users (
  id bigserial primary key,
  username varchar(255) not null,
  first_name varchar(225) not null,
  last_name varchar(225) not null,
  email varchar(255) not null,
  password_hash varchar(255) not null,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz
);

create unique index users_email_uniq_idx ON users(LOWER(email));
create unique index users_username_uniq_idx ON users(LOWER(username));

create index users_email_idx ON users(LOWER(email));
create index users_username_idx ON users(LOWER(username));

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop index if exists users_username_idx;
drop index if exists users_email_idx;

drop index if exists users_username_uniq_idx;
drop index if exists users_email_uniq_idx;

drop table if exists users;
