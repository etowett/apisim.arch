-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table messages (
  id bigserial primary key,
  user_id bigint not null references users(id),
  sender_id varchar(20) not null,
  meta varchar(999) not null,
  message varchar(999) not null,
  recipient_count integer not null,
  cost numeric(16,4) not null,
  currency varchar(20) not null,
  sent_at timestamptz not null,
  created_at timestamptz not null default clock_timestamp()
);

create index messages_user_id_idx ON messages(user_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop index if exists messages_user_id_idx;

drop table if exists messages;
