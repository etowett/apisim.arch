-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table sessions (
  id bigserial primary key,
  user_id bigint references users(id),
  deactivated_at timestamptz null,
  expires_at timestamptz not null,
  ip_address varchar(255) not null,
  last_refreshed_at timestamptz  not null,
  user_agent varchar(255) not null,
  created_at timestamptz  not null default clock_timestamp(),
  updated_at timestamptz  not null default clock_timestamp()
);

create index sessions_user_idx ON sessions(user_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop index if exists sessions_user_idx;

drop table if exists sessions;
