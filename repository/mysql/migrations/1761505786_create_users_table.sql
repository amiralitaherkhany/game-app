-- +migrate Up
-- Create users table
create table users(
                      id int primary key auto_increment,
                      name varchar(255) not null,
                      phone_number varchar(255) not null unique,
                      created_at timestamp default current_timestamp
);

-- +migrate Down
-- Drop users table
DROP TABLE users;