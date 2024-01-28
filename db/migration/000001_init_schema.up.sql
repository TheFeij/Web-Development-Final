create table users (
    id bigint primary key,
    firstname varchar(50),
    lastname varchar(50),
    phone varchar(11) unique,
    username varchar(128) unique,
    password varchar,
    image varchar unique,
    bio varchar(100)
);

create table contacts (
    user_id bigint,
    contact_id bigint,
    contact_name varchar(100),
    display_number boolean default true,
    display_profile_picture boolean default true,
    primary key (user_id, contact_id),
    foreign key (user_id) references users(id),
    foreign key (contact_id) references users(id)
);

create table chats (
    id bigint primary key,
    createdAt date default now()
);

create table chat_participants (
    chat_id bigint references chats(id),
    user_id bigint references users(id),
    primary key (chat_id, user_id)
);

create table chat_messages (
    id bigint primary key,
    chat_id bigint references chats(id),
    sender_id bigint,
    receiver_id bigint,
    content varchar(300),
    createdAt date default now(),
    foreign key (chat_id, sender_id) references chat_participants(chat_id, user_id),
    foreign key (chat_id, receiver_id) references chat_participants(chat_id, user_id)
);

create table groups (
    id bigint primary key,
    owner_id bigint unique,
    createdAt date default now()
);

create table group_participants (
    group_id bigint references groups(id),
    user_id bigint references users(id),
    primary key (group_id, user_id)
);

-- adding the foreign key that owner of a chat be a participant of that group
alter table groups
    add constraint fk_owner foreign key(id, owner_id) references group_participants(group_id, user_id);

create table group_messages (
    id bigint primary key,
    group_id bigint,
    sender_id bigint,
    content varchar(300),
    createdAt date default now(),
    foreign key (group_id, sender_id) references group_participants(group_id, user_id)
);

create table channels (
    id bigint primary key,
    owner_id bigint references users(id),
    createdAt date default now()
);

create table channel_participants (
    channel_id bigint references channels(id),
    user_id bigint references users(id),
    primary key (channel_id, user_id)
);

create table channel_admins (
    channel_id bigint references channels(id),
    user_id bigint references users(id),
    primary key (channel_id, user_id),
    foreign key (channel_id, user_id) references channel_participants(channel_id, user_id)
);

-- adding the foreign key that owner of a channel be an admin of that channel
alter table channels
    add constraint fk_owner foreign key(id, owner_id) references channel_admins(channel_id, user_id);

create table channel_posts (
    id bigint primary key,
    channel_id bigint references channels(id),
    sender_id bigint references users(id),
    content varchar(300),
    createdAt date default now(),
    foreign key (channel_id, sender_id) references channel_admins(channel_id, user_id)
);