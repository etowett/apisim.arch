-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table api_keys(
  id bigserial primary key,
  user_id bigint not null references users(id),
  provider varchar(100) not null,
  name varchar(100) not null,
  access_id varchar(100) not null,
  access_secret_hash varchar(256) not null,
  dlr_url varchar(255),
  deleted boolean default false,
  created_at timestamptz not null default clock_timestamp()
);

create index  api_keys_user_idx on api_keys(user_id);
create unique index api_keys_name_uniq_idx ON api_keys(user_id, lower(name));
create unique index api_keys_access_id_uniq_idx ON api_keys(user_id, provider, lower(access_id));

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop index if exists api_keys_access_id_uniq_idx;
drop index if exists api_keys_name_uniq_idx;
drop index if exists api_keys_user_idx;

drop table if exists api_keys;
