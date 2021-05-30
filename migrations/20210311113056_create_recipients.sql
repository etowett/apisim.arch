-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table recipients (
  id bigserial primary key,
  message_id bigint not null references messages (id),
  phone varchar(20) not null,
  api_id varchar(255) not null,
  route varchar(10) not null,
  cost numeric(16,4) not null,
  currency varchar(20) not null,
  created_at timestamptz not null default clock_timestamp()
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists recipients;
