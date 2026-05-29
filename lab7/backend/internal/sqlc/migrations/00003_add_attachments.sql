-- +goose Up
-- +goose StatementBegin
create table attachments
(
    id              bigserial primary key,
    filename        varchar(255)    not null,
    type_id         bigint          not null,
    created_at      timestamp       not null default now(),
    updated_at      timestamp       not null default now(),
    is_deleted      bool            not null default false,
    constraint fk_attachment_type foreign key (type_id) references attachment_types (id) on delete restrict
);
create unique index attachments_uk
    on attachments (filename)
    where attachments.is_deleted = false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table attachments;
-- +goose StatementEnd
