drop table if exists contacts cascade;

drop table if exists chat_messages cascade;

drop table if exists chat_participants cascade;

drop table if exists chats cascade;

alter table groups
    drop constraint if exists fk_owner cascade;

drop table if exists group_messages cascade;

drop table if exists group_participants cascade;

drop table if exists groups cascade;

alter table channels
    drop constraint if exists fk_owner cascade;

drop table if exists channel_posts cascade;

drop table if exists channel_admins cascade;

drop table if exists channel_participants cascade;

drop table if exists channels cascade;

drop table if exists users cascade;

