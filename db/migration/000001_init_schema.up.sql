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
    owner_id bigint references users(id),
    people bigint[] references users(id),
    createdAt date default now(),
    constraint people_length check (array_length(people, 1) = 1024)
);

create table group_messages (
    id bigint primary key,
    group_id bigint references groups(id),
    sender_id bigint references users(id),
    content varchar(300),
    createdAt date default now(),
    foreign key (group_id, sender_id) references groups(id, people)
);

create table channels (
    id bigint primary key,
    owner_id bigint references users(id),
    admins bigint[] references users(id),
    people bigint[] references users(id),
    createdAt date default now(),
    constraint people_length check (array_length(people, 1) = 1024),
    constraint admin_length check (array_length(admins, 1) = 32)
);

create table channel_messages (
    id bigint primary key,
    channel_id bigint references channels(id),
    sender_id bigint references users(id),
    content varchar(300),
    createdAt date default now(),
    foreign key (channel_id, sender_id) references channels(id, admins)
);