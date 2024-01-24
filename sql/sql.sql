CREATE DATABASE IF NOT EXISTS super_chat;
USE super_chat;

DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id int auto_increment primary key,
    email varchar(50) not null unique,
    password varchar(100) not null
);