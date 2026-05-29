-- +goose Up
-- +goose StatementBegin
create table posts
(
    id              bigserial primary key,
    title           varchar(255) not null,
    content         text         not null,
    creator_id      bigint       not null,
    created_at      timestamp    not null default now(),
    updated_at      timestamp    not null default now(),
    is_deleted      bool         not null default false,

    constraint fk_post_creator foreign key (creator_id) references users (id) on delete restrict
);

create table post_attachments
(
    post_id         bigint not null,
    attachment_id   bigint not null,

    constraint fk_attached_post foreign key (post_id) references posts (id) on delete restrict,
    constraint fk_post_attachment foreign key (attachment_id) references attachments (id) on delete restrict
);

create unique index post_attachment_uk
    on post_attachments (post_id, attachment_id);

create table tags
(
    id              bigserial primary key,
    name            varchar(255)    not null,
    created_at      timestamp       not null default now(),
    updated_at      timestamp       not null default now(),
    is_deleted      bool            not null default false
);

create unique index tags_uk
    on tags (name)
    where tags.is_deleted = false;

create table post_tags
(
    post_id     bigint not null,
    tag_id      bigint not null,

    constraint fk_tagged_post foreign key (post_id) references posts (id) on delete restrict,
    constraint fk_tag_id foreign key (tag_id) references tags (id) on delete restrict
);

create unique index post_tags_uk
    on post_tags (post_id, tag_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table post_attachments;
drop table post_tags;
drop table posts;
drop table tags;
-- +goose StatementEnd
