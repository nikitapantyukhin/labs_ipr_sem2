-- +goose Up
-- +goose StatementBegin
create table workouts
(
    id              bigserial primary key,
    club_id         bigint      not null,
    start_date      timestamp   not null,
    end_date        timestamp   not null,
    cancelled       bool        not null default false,
    created_at      timestamp   not null default now(),
    updated_at      timestamp   not null default now(),
    is_deleted      bool        not null default false,

    constraint fk_workout_club_id foreign key (club_id) references clubs (id) on delete restrict
);
create table metrics
(
    id                  bigserial primary key,
    name                varchar(255)    not null,
    description         text            not null,
    units               varchar(32)     not null,
    club_id             bigint          not null,
    created_at          timestamp       not null default now(),
    updated_at          timestamp       not null default now(),
    is_deleted          bool            not null default false,

    constraint fk_metrics_club foreign key (club_id) references clubs (id) on delete restrict
);

create unique index metrics_uk
    on metrics (name, description, units, club_id)
    where metrics.is_deleted = false;

create table workout_attendees
(
    id              bigserial primary key,
    workout_id      bigint      not null,
    request_id      bigint      not null,
    visited         bool        not null default false,
    review          text,
    created_at      timestamp   not null default now(),
    updated_at      timestamp   not null default now(),
    is_deleted      bool        not null default false,

    constraint fk_attended_workout_id foreign key (workout_id) references workouts (id) on delete restrict,
    constraint fk_club_request_id foreign key (request_id) references club_join_requests (id) on delete restrict
);

create unique index workout_attendees_uk
    on workout_attendees (workout_id, request_id)
    where workout_attendees.is_deleted = false;
create table workout_metrics
(
    workout_attendee_id         bigint                    not null,
    metric_id                   bigint                    not null,
    value                       double precision          not null,
    created_at                  timestamp                 not null default now(),
    updated_at                  timestamp                 not null default now(),
    is_deleted                  bool                      not null default false,

    constraint fk_attendee_id foreign key (workout_attendee_id) references workout_attendees (id) on delete restrict,
    constraint fk_workout_metric_id foreign key (metric_id) references metrics (id) on delete restrict
);

create unique index workout_metrics_pk
    on workout_metrics (workout_attendee_id, metric_id)
    where workout_metrics.is_deleted = false;

create table metrics_targets
(
    metric_id       bigint                   not null,
    user_id         bigint                   not null,
    value           double precision         not null,
    created_at      timestamp                not null default now(),
    updated_at      timestamp                not null default now(),
    is_deleted      bool                     not null default false,

    constraint fk_target_metric_id foreign key (metric_id) references metrics (id) on delete restrict,
    constraint fk_target_user_id foreign key (user_id) references users (id) on delete restrict
);

create unique index metrics_targets_uk
    on metrics_targets (metric_id, user_id)
    where metrics_targets.is_deleted = false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table workout_metrics;
drop table metrics_targets;
drop table workout_attendees;
drop table metrics;
drop table workouts;
-- +goose StatementEnd
