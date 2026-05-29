-- +goose Up
-- +goose StatementBegin
create table event_types
(
    id              bigserial primary key,
    name            varchar(255)    not null,
    created_at      timestamp       not null default now(),
    updated_at      timestamp       not null default now(),
    is_deleted      bool            not null default false
);

create unique index event_types_uk
    on event_types (name)
    where event_types.is_deleted = false;

create table events
(
    id                  bigserial primary key,
    name                varchar(255)    not null,
    description         text            not null,
    start_date          timestamp       not null,
    end_date            timestamp       not null,
    livestream_link     text,
    sport_type_id       bigint          not null,
    creator_id          bigint          not null,
    type_id             bigint          not null,
    total_places        int,
    place               text            not null,
    created_at          timestamp       not null default now(),
    updated_at          timestamp       not null default now(),
    is_deleted          bool            not null default false,

    constraint fk_event_sport_type foreign key (sport_type_id) references sport_types (id) on delete restrict,
    constraint fk_event_creator foreign key (creator_id) references users (id) on delete restrict,
    constraint fk_event_type_id foreign key (type_id) references event_types (id) on delete restrict
);

create table achievements
(
    id                  bigserial primary key,
    user_id             bigint      not null,
    event_id            bigint      not null,
    place               int         not null,
    created_at          timestamp   not null default now(),
    updated_at          timestamp   not null default now(),
    is_deleted          bool        not null default false,

    constraint fk_event_achievement foreign key (event_id) references events (id) on delete restrict,
    constraint fk_user_achievement foreign key (user_id) references users (id) on delete restrict
);

create unique index achievements_uk
    on achievements (event_id, user_id, place)
    where achievements.is_deleted = false;

create table event_join_requests
(
    id          bigserial primary key, 
    event_id    bigint       not null,
    user_id     bigint       not null,
    status      varchar(128) not null,
    created_at  timestamp    not null default now(),
    updated_at  timestamp    not null default now(),
    is_deleted  bool         not null default false,

    constraint fk_joined_event foreign key (event_id) references events (id) on delete restrict,
    constraint fk_user_joining_event foreign key (user_id) references users (id) on delete restrict
);

create unique index event_join_requests_uk
    on event_join_requests (event_id, user_id)
    where event_join_requests.is_deleted = false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table event_join_requests;
drop table achievements;
drop table events;
drop table event_types;
-- +goose StatementEnd
