create table users (
    id bigint primary key,
    firstname varchar(50),
    lastname varchar(50),
    phone char(11) unique,
    username varchar(128) unique,
    password varchar,
    image varchar,
    bio varchar(100),
    display_number boolean default true,
    display_profile_picture boolean default true
);

create table contacts (
    user_id bigint,
    contact_id bigint,
    contact_name varchar(100),
    primary key (user_id, contact_id),
    foreign key (user_id) references users(id),
    foreign key (contact_id) references users(id)
);

create table chats (
    id bigint primary key,
    is_dead boolean default false,
    created_at date default now()
);

create table chat_participants (
    chat_id bigint references chats(id),
    user_id bigint references users(id),
    primary key (chat_id, user_id)
);

create table chat_messages (
    id bigint primary key,
    chat_id bigint references chats(id),
    source_sender_id bigint references users(id),
    original_message_id bigint,
    sender_id bigint,
    receiver_id bigint,
    content varchar(300),
    created_at date default now(),
    foreign key (chat_id, sender_id) references chat_participants(chat_id, user_id),
    foreign key (chat_id, receiver_id) references chat_participants(chat_id, user_id)
);

create table groups (
    id bigint primary key,
    name varchar(64) not null,
    owner_id bigint references users(id),
    created_at date default now()
);

create table group_participants (
    group_id bigint references groups(id),
    user_id bigint references users(id),
    primary key (group_id, user_id)
);

create table group_messages (
    id bigint primary key,
    group_id bigint,
    source_sender_id bigint references users(id),
    original_message_id bigint,
    sender_id bigint,
    content varchar(300),
    created_at date default now(),
    foreign key (group_id, sender_id) references group_participants(group_id, user_id)
);

create table channels (
    id bigint primary key,
    name varchar(64) not null,
    owner_id bigint references users(id),
    created_at date default now()
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

create table channel_posts (
    id bigint primary key,
    channel_id bigint references channels(id),
    source_sender_id bigint references users(id),
    original_message_id bigint,
    sender_id bigint references users(id),
    content varchar(300),
    created_at date default now(),
    foreign key (channel_id, sender_id) references channel_admins(channel_id, user_id)
);