-- +goose Up
-- +goose StatementBegin
create table club_join_requests
(
    id         bigserial primary key, 
    club_id    bigint       not null,
    user_id    bigint       not null,
    status     varchar(128) not null,
    created_at timestamp    not null default now(),
    updated_at timestamp    not null default now(),
    is_deleted bool         not null default false,

    constraint fk_user_joining_club foreign key (user_id) references users (id) on delete restrict,
    constraint fk_joined_club foreign key (club_id) references clubs (id) on delete restrict
);
create unique index club_join_requests_uk
    on club_join_requests (club_id, user_id)
    where club_join_requests.is_deleted = false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table club_join_requests;
-- +goose StatementEnd
