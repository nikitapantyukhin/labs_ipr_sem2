-- +goose Up
-- +goose StatementBegin
create table group_types
(
    id         bigserial primary key,
    name       varchar(128) not null,
    created_at timestamp    not null default now(),
    updated_at timestamp    not null default now(),
    is_deleted bool         not null default false
);

create unique index group_types_name_uk
    on group_types (name)
    where group_types.is_deleted = false;

create table groups
(
    id              bigserial primary key,
    institute       varchar(64) not null,
    enrollment_year integer     not null,
    prefix          varchar(32),
    group_type_id   bigint      not null,
    group_number    smallint,
    created_at      timestamp   not null default now(),
    updated_at      timestamp   not null default now(),
    is_deleted      bool        not null default false,

    constraint fk_group_type foreign key (group_type_id) references group_types (id) on delete restrict
);

create unique index groups_uk
    on groups (institute, enrollment_year, prefix, group_type_id, group_number)
    where groups.is_deleted = false;

create table users
(
    id                  bigserial primary key,
    full_name           varchar(255) not null,
    social_network_link varchar(255) not null,
    phone_number        varchar(32)  not null,
    email               varchar(255) not null,
    birth_date          timestamp    not null,
    role                varchar(128) not null,
    password            bytea        not null,
    group_id            bigint,
    created_at          timestamp    not null default now(),
    updated_at          timestamp    not null default now(),
    is_deleted          bool         not null default false,

    constraint fk_group foreign key (group_id) references groups (id) on delete restrict
);

create unique index users_email_uk
    on users (email)
    where users.is_deleted = false;

create unique index users_social_network_link_uk
    on users (social_network_link)
    where users.is_deleted = false;

create unique index users_phone_number_uk
    on users (phone_number)
    where users.is_deleted = false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop table groups;
drop table group_types;
-- +goose StatementEnd
