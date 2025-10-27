-- +migrate Up
-- add password column to users table
alter table users ADD COLUMN password varchar(255) not null;

-- +migrate Down
-- Drop password column from users table
alter table users DROP column password;