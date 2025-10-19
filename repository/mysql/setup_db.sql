create table users(
    id int primary key auto_increment,
    name varchar(255) not null,
    phone_number varchar(255) not null unique,
    password varchar(255) not null ,
    created_at timestamp default current_timestamp
);