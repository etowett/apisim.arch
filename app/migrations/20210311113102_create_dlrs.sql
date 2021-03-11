-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table dlrs (
  id bigserial primary key,
  recipient_id bigint not null references recipients (id),
  status varchar(20) not null,
  reason varchar(255) not null,
  received_at timestamptz not null,
  created_at timestamptz not null default clock_timestamp()
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists dlrs;
