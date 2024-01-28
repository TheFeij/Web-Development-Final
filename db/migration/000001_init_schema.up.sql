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