create table users (
    id bigint primary key,
    firstname varchar(50),
    lastname varchar(50),
    phone varchar(11) unique,
    username varchar(128) unique,
    password varchar,
    image varchar unique,
    bio varchar(100)
)