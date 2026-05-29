-- +goose Up
-- +goose StatementBegin
create table club_attachments
(
    club_id         bigint not null,
    attachment_url   varchar(255) not null,

    constraint fk_attached_club foreign key (club_id) references clubs (id) on delete restrict
);

create unique index club_attachment_uk
    on club_attachments (club_id, attachment_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table club_attachments;
-- +goose StatementEnd
