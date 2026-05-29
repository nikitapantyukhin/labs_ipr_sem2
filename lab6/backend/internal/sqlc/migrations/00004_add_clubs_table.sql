-- +goose Up
-- +goose StatementBegin
create table sport_types
(
    id         bigserial primary key, 
    name       varchar(255) not null,
    created_at timestamp    not null default now(),
    updated_at timestamp    not null default now(),
    is_deleted bool         not null default false
);

create unique index sport_types_uk
    on sport_types (name)
    where sport_types.is_deleted = false;

create table education_levels
(
    id         bigserial primary key,
    name       varchar(255)  not null,
    created_at timestamp     not null default now(),
    updated_at timestamp     not null default now(),
    is_deleted bool          not null default false
);

create unique index education_levels_uk
    on education_levels (name)
    where education_levels.is_deleted = false;

create table clubs
(
    id                          bigserial primary key,
    name                        varchar(255) not null,
    description                 text         not null,
    sport_type_id               bigint       not null,
    teacher_id                  bigint       not null,
    total_places                int,
    place                       text         not null,
    education_level_id          bigint       not null,
    required_workout_per_week   int not      null,
    created_at                  timestamp    not null default now(),
    updated_at                  timestamp    not null default now(),
    is_deleted                  bool         not null default false,

    constraint fk_club_sport_type foreign key (sport_type_id) references sport_types (id) on delete restrict,
    constraint fk_teacher foreign key (teacher_id) references users (id) on delete restrict,
    constraint fk_education_level foreign key (education_level_id) references education_levels (id) on delete restrict
);

create unique index clubs_uk
    on clubs (name, sport_type_id, teacher_id, education_level_id)
    where clubs.is_deleted = false;

create table reviews
(
    id              bigserial primary key,
    rating          int       not null,
    content         text      not null,
    creator_id      bigint    not null,
    club_id         bigint    not null,
    created_at      timestamp not null default now(),
    updated_at      timestamp not null default now(),
    is_deleted      bool      not null default false,

    constraint fk_reviewer foreign key (creator_id) references users (id) on delete restrict,
    constraint fk_reviewed_club foreign key (club_id) references clubs (id) on delete restrict

);

create unique index reviews_uk
    on reviews (creator_id, club_id)
    where reviews.is_deleted = false;

create table review_attachments
(
    review_id       bigint not null,
    attachment_id   bigint not null,

    constraint fk_review_id foreign key (review_id) references reviews (id) on delete restrict,
    constraint fk_review_attachment_id foreign key (attachment_id) references attachments (id) on delete restrict
);
create unique index review_attachments_uk
    on review_attachments (review_id, attachment_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table review_attachments;
drop table reviews;
drop table clubs;
drop table education_levels;
drop table sport_types;
-- +goose StatementEnd
