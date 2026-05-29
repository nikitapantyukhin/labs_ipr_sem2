-- +goose Up
-- +goose StatementBegin
create table attachment_types
(
    id              bigserial primary key,
    name            varchar(128)    not null,
    created_at      timestamp       not null default now(),
    is_deleted      bool            not null default false
);

create unique index attachment_types_uk
    on attachment_types (name)
    where attachment_types.is_deleted = false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table attachment_types;
-- +goose StatementEnd
