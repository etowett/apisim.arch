-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd


create table transactions (
  id bigserial primary key,
  user_id bigint not null references users(id),
  amount numeric(16, 4) not null,
  balance numeric(16, 4) not null,
  currency varchar(10) not null,
  code varchar(40) not null,
  type varchar(40) not null,
  created_at timestamptz not null default clock_timestamp()
);

create index transactions_user_id_idx on transactions(user_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop index if exists transactions_user_id_idx;

drop table if exists transactions;
